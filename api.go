package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wolv89/chirpy/internal/database"
)

type Chirp struct {
	Body string `json:"body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func responseJSON(w http.ResponseWriter, status int, data interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	dat, jerr := json.Marshal(data)
	if jerr != nil {
		log.Printf("Error marshaling JSON %s", jerr)
		return
	}

	w.Write(dat)

}

func (cfg *apiConfig) APIHealthCheck(w http.ResponseWriter, _ *http.Request) {

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func (cfg *apiConfig) APICreateChirp(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	newChirp := database.CreateChirpParams{}
	err := decoder.Decode(&newChirp)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Something went wrong"})
		return
	}

	if len(newChirp.Body) == 0 {
		responseJSON(w, http.StatusBadRequest, ErrorResponse{"Silent chirp"})
		return
	}

	if len(newChirp.Body) > 140 {
		responseJSON(w, http.StatusBadRequest, ErrorResponse{"Chirp is too long"})
		return
	}

	newChirp.Body = FilterChirp(newChirp.Body)

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), newChirp)

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Unable to post chirp"})
		return
	}

	responseJSON(w, http.StatusCreated, chirp)

}

func (cfg *apiConfig) APICreateUser(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	type NewUser struct {
		Email string `json:"email"`
	}
	newUser := NewUser{}
	err := decoder.Decode(&newUser)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Something went wrong"})
		return
	}

	if len(newUser.Email) == 0 {
		responseJSON(w, http.StatusBadRequest, ErrorResponse{"New users must have an email address"})
		return
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), newUser.Email)

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Unable to create user"})
		return
	}

	responseJSON(w, http.StatusCreated, user)

}
