package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

var profaneWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

type chirpPost struct {
	Body string `json:"body"`
}

type serverError struct {
	E string `json:"error"`
}

type chirpValid struct {
	CleanedBody string `json:"cleaned_body"`
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	chirp := chirpPost{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respBody := serverError{
			E: "Something went wrong",
		}
		respondWithJson(w, http.StatusInternalServerError, respBody)
		return
	}

	if len(chirp.Body) > 140 {
		respBody := serverError{
			E: "Chirp is too long",
		}
		respondWithJson(w, http.StatusBadRequest, respBody)
		return
	}

	cleanChirp := filterProfanity(chirp.Body)
	respBody := chirpValid{
		CleanedBody: cleanChirp,
	}
	respondWithJson(w, http.StatusOK, respBody)
}

func filterProfanity(chirp string) string {
	chirpList := strings.Split(chirp, " ")
	for i := range len(chirpList) {
		loweredWord := strings.ToLower(chirpList[i])
		if _, ok := profaneWords[loweredWord]; ok {
			chirpList[i] = "****"
		}
	}
	return strings.Join(chirpList, " ")
}
