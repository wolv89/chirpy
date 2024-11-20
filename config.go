package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/wolv89/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetricsView(w http.ResponseWriter, _ *http.Request) {

	w.Header().Add("Content-Type", "text/html")

	fmt.Fprintf(w, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())

}

func (cfg *apiConfig) handleMetricsReset(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	fmt.Fprintf(w, "Reset")
}
