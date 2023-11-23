package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type userResource struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *ApiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	params, err := getParamsFromRequest(userResource{}, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = cfg.validateParams(w, params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := cfg.Db.CreateUser(params.Password, params.Email)
	if err != nil {
		message := fmt.Sprintf("Couldn't create user: %v", err)
		respondWithError(w, http.StatusInternalServerError, message)
		return
	}

	respondWithJSON(w, 201, user)
}

func (cfg *ApiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	user_id, err := cfg.stripTokenFromHeader(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := getParamsFromRequest(userResource{}, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.Db.UpdateUser(user_id, params.Password, params.Email)
	if err != nil {
		message := fmt.Sprintf("Couldn't create user: %v", err)
		respondWithError(w, http.StatusInternalServerError, message)
		return
	}

	respondWithJSON(w, 200, user)
}

func (cfg *ApiConfig) validateParams(w http.ResponseWriter, params userResource) error {
	if params.Password == "" {
		return errors.New("password cant be empty")
	}
	if params.Email == "" {
		return errors.New("email cant be empty")
	}

	_, err := cfg.Db.GetUserByEmail(params.Email)
	if err == nil {
		return errors.New("email is already in use")
	}

	return nil
}

func (cfg *ApiConfig) stripTokenFromHeader(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization") //Authorization: Bearer <token>

	auth_token := strings.Split(auth, " ")[1]

	token, err := jwt.ParseWithClaims(auth_token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT), nil
	})
	if err != nil {
		fmt.Println(err)
		return "", errors.New("user is not authorized")
	}

	user_id, err := token.Claims.GetSubject()
	if err != nil {
		return "", errors.New("couldn't get id from token")
	}
	return user_id, nil
}
