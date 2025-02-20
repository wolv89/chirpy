package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/wolv89/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits                atomic.Int32
	dbQueries                     *database.Queries
	platform, jwtSecret, polkaKey string
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

func (cfg *apiConfig) handleFullReset(w http.ResponseWriter, req *http.Request) {

	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	cfg.fileserverHits.Store(0)
	_ = cfg.dbQueries.DeleteUsers(req.Context())
	fmt.Fprintf(w, "Reset")
}
