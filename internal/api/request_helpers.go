package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	db "github.com/SilverLuhtoja/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits int
	Db             *db.DB
	JWT            string
	PolkaKey       string
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	respondWithJSON(w, code, map[string]string{"error": fmt.Sprintf("%+v", err)})
}

func getParamsFromRequest[T interface{}](structBody T, r *http.Request) (T, error) {
	decoder := json.NewDecoder(r.Body)

	params := structBody
	err := decoder.Decode(&params)
	if err != nil {
		return structBody, errors.New("couldn't decode parameters")
	}

	return params, nil
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

func serverReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
