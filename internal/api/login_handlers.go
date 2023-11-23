package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type loginResource struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type loginResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *ApiConfig) logIn(w http.ResponseWriter, r *http.Request) {
	params, err := getParamsFromRequest(loginResource{}, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.Db.GetUserByEmail(params.Email)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password)) != nil {
		respondWithError(w, http.StatusUnauthorized, "Email or password is wrong")
		return
	}

	full_day := 60 * 60 * 24
	if params.ExpiresInSeconds > full_day || params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = full_day
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(params.ExpiresInSeconds))),
		Subject:   fmt.Sprint(user.Id),
	})

	jwt, err := newToken.SignedString([]byte(cfg.JWT))
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Error with signing jwt")
		return
	}
	response := loginResponse{
		Id:    fmt.Sprint(user.Id),
		Email: user.Email,
		Token: jwt,
	}
	respondWithJSON(w, 200, response)
}
