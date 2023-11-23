package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	db "github.com/SilverLuhtoja/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits int
	Db             *db.DB
	JWT            string
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
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

func serverReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
