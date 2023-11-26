package api

import (
	"errors"
	"net/http"
	"strings"
)

type membershipResource struct {
	Event string         `json:"event"`
	Data  map[string]int `json:"data"`
}

func (cfg *ApiConfig) GrantMemberhip(w http.ResponseWriter, r *http.Request) {
	err := cfg.validateKey(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}

	params, err := getParamsFromRequest(membershipResource{}, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 200, "")
		return
	}

	err = cfg.Db.UpdateToMember(params.Data["user_id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, 200, "")
}

func (cfg *ApiConfig) validateKey(r *http.Request) error {
	auth_params := r.Header.Get("Authorization")
	if auth_params == "" {
		return errors.New("problematic authorization header")
	}

	req_key := strings.Split(auth_params, " ")[1]
	if req_key != cfg.PolkaKey {
		return errors.New("request denied")
	}

	return nil
}
