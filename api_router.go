package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	db "github.com/SilverLuhtoja/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

type message struct {
	Body string `json:"body"`
}

type apiConfig struct {
	fileserverHits int
	db             *db.DB
}

func newApiRouter(apiConfig *apiConfig) *chi.Mux {
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", serverReadiness)
	apiRouter.Get("/metrics", apiConfig.showMetrics)
	apiRouter.HandleFunc("/reset", apiConfig.resetMetrics)
	apiRouter.Post("/chirps", apiConfig.saveChirp)
	apiRouter.Get("/chirps", apiConfig.getChirps)
	return apiRouter
}

func serverReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) saveChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := message{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	//  if all good
	chirp, err := cfg.db.CreateChirp(cleanInput(params.Body))
	if err != nil {
		message := fmt.Sprintf("Couldn't create chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, message)
		return
	}
	respondWithJSON(w, 201, chirp)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	res, err := cfg.db.GetChirps()
	if err != nil {
		message := fmt.Sprintf("Couldn't get chirps: %v", err)
		respondWithError(w, http.StatusInternalServerError, message)
		return
	}
	respondWithJSON(w, 200, res)
}

func cleanInput(paramsBody string) string {
	message := strings.Split(paramsBody, " ")
	cleaned := []string{}
	bad_words := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range message {
		if slices.Contains(bad_words, strings.ToLower(word)) {
			cleaned = append(cleaned, "****")
		} else {
			cleaned = append(cleaned, word)
		}
	}
	return strings.Join(cleaned, " ")
}
