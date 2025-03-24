package main

import (
  _ "github.com/lib/pq"
	"log"
	"net/http"
	"sync/atomic"
  "os"
  "github.com/codybstrange/diy-server/internal/database"
  "database/sql"
  "github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
  db *database.Queries
  platform string
}

func main() {
	godotenv.Load()
  const filepathRoot = "."
	const port = "8080"

  dbURL := os.Getenv("DB_URL")
  if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

  db, err := sql.Open("postgres", dbURL)
  if err != nil {
    log.Fatal("Error in opening postgres connection")
  }
  dbQueries := database.New(db)

  platform := os.Getenv("PLATFORM")
  
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
    db: dbQueries,
    platform: platform,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
  mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
  mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerGetChirp)

  mux.HandleFunc("POST /api/chirps", apiCfg.handlerPostChirp)
  mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
