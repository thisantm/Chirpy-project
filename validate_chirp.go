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

type chirpValid struct {
	CleanedBody string `json:"cleaned_body"`
}

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	chirp := chirpPost{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusInternalServerError, "Chirp is too long", err)
		return
	}

	cleanChirp := filterProfanity(chirp.Body)
	respBody := chirpValid{
		CleanedBody: cleanChirp,
	}
	respondWithJSON(w, http.StatusOK, respBody)
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
