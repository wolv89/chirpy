package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {

	port := "8080"
	filepathRoot := "."

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetricsView)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleMetricsReset)

	mux.HandleFunc("GET /api/healthz", APIHealthCheck)
	mux.HandleFunc("POST /api/validate_chirp", APIValidateChirp)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
