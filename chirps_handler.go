package main

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func (cfg *apiConfig) handlerGetChirpById(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			respBody := serverError{
				E: "Something went wrong",
			}
			respondWithJson(w, http.StatusInternalServerError, respBody)
			return
		}

		chirp, err := db.GetChirpById(id)
		if err != nil {
			if err.Error() == "not found" {
				respBody := serverError{
					E: "Chirp Not Found",
				}
				respondWithJson(w, http.StatusNotFound, respBody)
				return
			}
			respBody := serverError{
				E: "Something went wrong",
			}
			respondWithJson(w, http.StatusInternalServerError, respBody)
			return
		}

		respondWithJson(w, http.StatusOK, chirp)
	}
}

func (cfg *apiConfig) handlerGetChirps(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		chirps, err := db.GetChirps()

		if err != nil {
			respBody := serverError{
				E: "Something went wrong",
			}
			respondWithJson(w, http.StatusInternalServerError, respBody)
			return
		}

		respondWithJson(w, http.StatusOK, chirps)
	}
}

func (cfg *apiConfig) handlerCreateChirp(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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
		respBody, err := db.CreateChirp(cleanChirp)
		if err != nil {
			respBody := serverError{
				E: "Something went wrong",
			}
			respondWithJson(w, http.StatusInternalServerError, respBody)
			return
		}

		respondWithJson(w, http.StatusCreated, respBody)
	}
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
