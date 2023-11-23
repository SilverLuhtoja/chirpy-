package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SilverLuhtoja/chirpy/internal/api"
	db "github.com/SilverLuhtoja/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	port := "8080"
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	db := db.NewDB("database.json")
	apiConfig := &api.ApiConfig{Db: db, JWT: jwtSecret}

	router := chi.NewRouter()
	apiRouter := api.NewApiRouter(apiConfig)
	amdminRouter := chi.NewRouter()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	router.Handle("/app/*", apiConfig.MiddlewareMetricsInc(fsHandler))
	router.Handle("/app", apiConfig.MiddlewareMetricsInc(fsHandler))
	amdminRouter.Get("/metrics", apiConfig.RenderAdminMetrics)
	router.Mount("/api", apiRouter)
	router.Mount("/admin", amdminRouter)

	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr:        "localhost:" + port,
		Handler:     corsMux,
		ReadTimeout: 2 * time.Second,
	}

	log.Printf("Server running on: http://localhost:%s/app\n", port)
	log.Fatal(server.ListenAndServe())
}
