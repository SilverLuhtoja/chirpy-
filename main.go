package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/SilverLuhtoja/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := db.NewDB("database.json")
	if err != nil {
		log.Fatal("couldn't create database")
	}

	port := "8080"
	apiConfig := &apiConfig{db: db}
	router := chi.NewRouter()
	apiRouter := newApiRouter(apiConfig)
	amdminRouter := chi.NewRouter()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	router.Handle("/app/*", apiConfig.middlewareMetricsInc(fsHandler))
	router.Handle("/app", apiConfig.middlewareMetricsInc(fsHandler))
	amdminRouter.Get("/metrics", apiConfig.renderAdminMetrics)
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

func (cfg *apiConfig) renderAdminMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, cfg.fileserverHits)
	w.Write([]byte(message))
}
