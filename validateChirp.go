package main

import (
	"encoding/json"
	"net/http"
)

type chirpPost struct {
	Body string `json:"body"`
}

type serverError struct {
	E string `json:"error"`
}

type chirpValid struct {
	Valid bool `json:"valid"`
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

	respBody := chirpValid{
		Valid: true,
	}
	respondWithJson(w, http.StatusOK, respBody)
}
