package api

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type loginResource struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type loginResponse struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshResponse struct {
	Token string `json:"token"`
}

const ACCESS_ISSUER = "chirpy-access"
const REFRESH_ISSUER = "chirpy-refresh"

func (cfg *ApiConfig) logIn(w http.ResponseWriter, r *http.Request) {
	params, err := getParamsFromRequest(loginResource{}, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	user, err := cfg.Db.GetUserByEmail(params.Email)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	if bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password)) != nil {
		respondWithError(w, http.StatusUnauthorized, errors.New("email or password is wrong"))
		return
	}

	access_token, err := cfg.createAccessToken(ACCESS_ISSUER, user.Id)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}

	refresh_token, err := cfg.createAccessToken(REFRESH_ISSUER, user.Id)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}

	response := loginResponse{
		Id:           fmt.Sprint(user.Id),
		Email:        user.Email,
		Token:        access_token,
		RefreshToken: refresh_token,
	}
	respondWithJSON(w, 200, response)
}

func (cfg *ApiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	token, err := cfg.getTokenFromHeader(REFRESH_ISSUER, r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}

	tokens, err := cfg.Db.GetTokens()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	for _, tk := range tokens {
		if token.Raw == tk {
			respondWithError(w, http.StatusUnauthorized, errors.New("access denied"))
			return
		}
	}

	new_access_token, err := cfg.createAccessTokenFromHeader(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, 200, refreshResponse{Token: new_access_token})
}

func (cfg *ApiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	token, err := cfg.getTokenFromHeader(REFRESH_ISSUER, r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}

	err = cfg.Db.RevokeToken(token.Raw)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, 200, "")
}
