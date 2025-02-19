package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wolv89/chirpy/internal/database"
)

func main() {

	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	port := "8080"
	filepathRoot := "."

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
		platform:       os.Getenv("PLATFORM"),
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetricsView)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleFullReset)

	mux.HandleFunc("GET /api/healthz", apiCfg.APIHealthCheck)

	mux.HandleFunc("POST /api/login", apiCfg.APILogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.APIRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.APIRevoke)

	mux.HandleFunc("POST /api/users", apiCfg.APICreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.APIUpdateUser)

	mux.HandleFunc("GET /api/chirps", apiCfg.APIGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.APIGetChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.APICreateChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpId}", apiCfg.APIDeleteChirp)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
