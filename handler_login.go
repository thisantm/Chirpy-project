package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/thisantm/Chirpy-project/internal/auth"
)

type LoginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds string `json:"expires_in_seconds"`
}

type LoginResponse struct {
	User
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	loginRequest := LoginRequest{}
	err := decoder.Decode(&loginRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	timeToExpire, err := time.ParseDuration(loginRequest.ExpiresInSeconds)
	if err != nil || timeToExpire > time.Hour {
		timeToExpire = time.Hour
	}

	dbUser, err := cfg.db.GetUserByEmail(req.Context(), loginRequest.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user does not exist", err)
		return
	}

	err = auth.CheckPasswordHash(loginRequest.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "wrong password", err)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, timeToExpire)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to create login token", err)
		return
	}

	response := LoginResponse{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, response)
}
