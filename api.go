package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

func APIHealthCheck(w http.ResponseWriter, _ *http.Request) {

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func APIValidateChirp(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	chirp := Chirp{}
	err := decoder.Decode(&chirp)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Something went wrong"})
		return
	}

	if len(chirp.Body) == 0 {
		responseJSON(w, http.StatusBadRequest, ErrorResponse{"Silent chirp"})
		return
	}

	if len(chirp.Body) > 140 {
		responseJSON(w, http.StatusBadRequest, ErrorResponse{"Chirp is too long"})
		return
	}

	bannedWords := make(map[string]interface{})
	bannedWords["kerfuffle"] = nil
	bannedWords["sharbert"] = nil
	bannedWords["fornax"] = nil

	chirpWords := strings.Split(chirp.Body, " ")

	var checkWord string
	var ok bool

	for w := 0; w < len(chirpWords); w++ {
		checkWord = strings.ToLower(chirpWords[w])
		if _, ok = bannedWords[checkWord]; ok {
			chirpWords[w] = "****"
		}
	}

	responseJSON(w, http.StatusOK, ValidResponse{strings.Join(chirpWords, " ")})

}
