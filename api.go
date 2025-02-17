package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wolv89/chirpy/internal/auth"
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

type BasicAuth struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type AuthResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
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

func (cfg *apiConfig) getUserFromAuth(req *http.Request) (uuid.UUID, error) {

	blank := uuid.UUID{}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		return blank, err
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		return blank, err
	}

	return userId, nil

}

func (cfg *apiConfig) APILogin(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	userLogin := BasicAuth{}
	err := decoder.Decode(&userLogin)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Something went wrong"})
		return
	}

	if len(userLogin.Email) == 0 || len(userLogin.Password) == 0 {
		responseJSON(w, http.StatusBadRequest, ErrorResponse{"Please enter email and password"})
		return
	}

	comparePassword, err := cfg.dbQueries.GetUserPasswordFromEmail(req.Context(), userLogin.Email)

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Something went wrong"})
		return
	}

	loginError := auth.CheckPasswordHash(userLogin.Password, comparePassword)

	if loginError != nil {
		responseJSON(w, http.StatusUnauthorized, ErrorResponse{"incorrect email or password"})
		return
	}

	user, err := cfg.dbQueries.GetUserFromEmail(req.Context(), userLogin.Email)

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Unable to load user"})
		return
	}

	expiry := userLogin.ExpiresInSeconds
	if expiry <= 0 || expiry > 3600 {
		expiry = 3600
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Second*time.Duration(expiry))
	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Unable to generate auth token"})
		return
	}

	ar := AuthResponse{
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Email,
		token,
	}

	responseJSON(w, http.StatusOK, ar)

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

	userId, err := cfg.getUserFromAuth(req)
	if err != nil {
		fmt.Println(err)
		responseJSON(w, http.StatusUnauthorized, nil)
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
	newChirp.UserID = uuid.NullUUID{
		UUID:  userId,
		Valid: true,
	}

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

	newUser := BasicAuth{}
	err := decoder.Decode(&newUser)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Something went wrong"})
		return
	}

	if len(newUser.Email) == 0 || len(newUser.Password) == 0 {
		responseJSON(w, http.StatusBadRequest, ErrorResponse{"New users must have an email address and password"})
		return
	}

	password, err := auth.HashPassword(newUser.Password)

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Unusable password"})
		return
	}

	newUserParams := database.CreateUserParams{
		Email:          newUser.Email,
		HashedPassword: password,
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), newUserParams)

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Unable to create user"})
		return
	}

	responseJSON(w, http.StatusCreated, user)

}

func (cfg *apiConfig) APIGetAllChirps(w http.ResponseWriter, req *http.Request) {

	chirps, err := cfg.dbQueries.GetAllChirps(req.Context())

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	responseJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) APIGetChirp(w http.ResponseWriter, req *http.Request) {

	qry := req.PathValue("chirpId")

	if len(qry) == 0 {
		responseJSON(w, http.StatusBadRequest, ErrorResponse{"Need a chirp ID!"})
		return
	}

	uuid, err := uuid.Parse(qry)
	if err != nil {
		responseJSON(w, http.StatusBadRequest, ErrorResponse{"That doesn't look like a chirp ID!"})
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), uuid)

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	responseJSON(w, http.StatusOK, chirp)

}
