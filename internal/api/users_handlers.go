package api

import (
	"fmt"
	"net/http"
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
	user_id, err := cfg.validateToken(r)
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
