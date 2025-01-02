package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thisantm/Chirpy-project/internal/auth"
	"github.com/thisantm/Chirpy-project/internal/database"
)

var profaneWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type chirpPost struct {
	Body string `json:"body"`
}

type chirpResponse struct {
	Chirp
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	chirpPost := chirpPost{}
	err := decoder.Decode(&chirpPost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	bearer, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "authorization header not found", err)
		return
	}

	userId, err := auth.ValidateJWT(bearer, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid login token", err)
		return
	}

	chirpClean, err := handlerValidateChirp(chirpPost.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error validating chirp", err)
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   chirpClean,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	chirpResponse := chirpResponse{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	}

	respondWithJSON(w, http.StatusCreated, chirpResponse)
}

func handlerValidateChirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		return "", errors.New("Chirp is too long")
	}

	chirpClean := filterProfanity(chirp)
	return chirpClean, nil
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
