package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func APIHealthCheck(w http.ResponseWriter, _ *http.Request) {

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func APIValidateChirp(w http.ResponseWriter, req *http.Request) {

	type Chirp struct {
		Body string `json:"body"`
	}

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	type ValidResponse struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	chirp := Chirp{}
	err := decoder.Decode(&chirp)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errResp := ErrorResponse{"Something went wrong"}
		dat, jerr := json.Marshal(errResp)
		if jerr != nil {
			log.Printf("Error marshaling JSON %s", jerr)
			return
		}
		w.Write(dat)
		return
	}

	if len(chirp.Body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		errResp := ErrorResponse{"Silent chirp"}
		dat, jerr := json.Marshal(errResp)
		if jerr != nil {
			log.Printf("Error marshaling JSON %s", jerr)
			return
		}
		w.Write(dat)
		return
	}

	if len(chirp.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		errResp := ErrorResponse{"Chirp is too long"}
		dat, jerr := json.Marshal(errResp)
		if jerr != nil {
			log.Printf("Error marshaling JSON %s", jerr)
			return
		}
		w.Write(dat)
		return
	}

	w.WriteHeader(http.StatusOK)
	validResp := ValidResponse{true}
	dat, jerr := json.Marshal(validResp)
	if jerr != nil {
		log.Printf("Error marshaling JSON %s", jerr)
		return
	}
	w.Write(dat)

}
