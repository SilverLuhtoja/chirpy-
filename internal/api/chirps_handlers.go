package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type messageChirp struct {
	AuthorId string `json:"author_id"`
	Body     string `json:"body"`
}

func (cfg *ApiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.Db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	param, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	for _, chirp := range chirps {
		if chirp.Id == param {
			respondWithJSON(w, 200, chirp)
			return
		}
	}
	http.NotFound(w, r)
}

func (cfg *ApiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := messageChirp{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("couldn't decode parameters"))
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, errors.New("chirp is too long"))
		return
	}

	user_id, err := cfg.getUserIdFromAccessToken(r)
	if err != nil {
		respondWithError(w, 400, err)
		return
	}

	chirp, err := cfg.Db.SaveChirp(user_id, cleanInput(params.Body))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, 201, chirp)
}

func (cfg *ApiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	res, err := cfg.Db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, 200, res)
}

func (cfg *ApiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	user_id, err := cfg.getUserIdFromAccessToken(r)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err)
		return
	}

	param_id, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	chirp, err := cfg.Db.GetChirpById(param_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	if chirp.AuthorId != user_id {
		respondWithError(w, http.StatusForbidden, errors.New("no permissions to delete"))
		return
	}

	err = cfg.Db.DeleteChirp(param_id)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err)
		return
	}

	respondWithJSON(w, 200, "")
}
