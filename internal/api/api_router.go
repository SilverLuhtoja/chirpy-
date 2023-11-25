package api

import (
	"github.com/go-chi/chi/v5"
)

func NewApiRouter(apiConfig *ApiConfig) *chi.Mux {
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", serverReadiness)
	apiRouter.Get("/metrics", apiConfig.showMetrics)
	apiRouter.HandleFunc("/reset", apiConfig.resetMetrics)

	apiRouter.Post("/chirps", apiConfig.createChirp)
	apiRouter.Get("/chirps", apiConfig.getChirps)
	apiRouter.Get("/chirps/{chirpID}", apiConfig.getChirpById)
	apiRouter.Delete("/chirps/{chirpID}", apiConfig.deleteChirp)

	apiRouter.Post("/users", apiConfig.createUser)
	apiRouter.Put("/users", apiConfig.updateUser)
	apiRouter.Post("/login", apiConfig.logIn)
	apiRouter.Post("/refresh", apiConfig.refresh)
	apiRouter.Post("/revoke", apiConfig.revoke)

	return apiRouter
}
