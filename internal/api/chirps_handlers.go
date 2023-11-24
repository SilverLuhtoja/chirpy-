package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type messageChirp struct {
	Body string `json:"body"`
}

func (cfg *ApiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.Db.GetChirps()
	if err != nil {
		message := fmt.Sprintf("Couldn't get chirps: %v", err)
		respondWithError(w, http.StatusInternalServerError, message)
		return
	}

	param, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		message := fmt.Sprintf("error with conversion: %v", err)
		respondWithError(w, http.StatusInternalServerError, message)
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	chirp, err := cfg.Db.SaveChirp(cleanInput(params.Body))
	if err != nil {
		message := fmt.Sprintf("Couldn't create chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, message)
		return
	}
	respondWithJSON(w, 201, chirp)
}

func (cfg *ApiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	res, err := cfg.Db.GetChirps()
	if err != nil {
		message := fmt.Sprintf("Couldn't get chirps: %v", err)
		respondWithError(w, http.StatusInternalServerError, message)
		return
	}
	respondWithJSON(w, 200, res)
}
