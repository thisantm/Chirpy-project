package main

import (
	"encoding/json"
	"net/http"

	"github.com/thisantm/Chirpy-project/internal/auth"
	"github.com/thisantm/Chirpy-project/internal/database"
)

func (cfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, req *http.Request) {
	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed getting access token", err)
		return
	}

	userId, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid access token", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	userRequest := UserRequest{}
	err = decoder.Decode(&userRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to hash password", err)
		return
	}

	updateUserEmailAndPasswordParams := database.UpdateUserEmailAndPasswordParams{
		Email:          userRequest.Email,
		HashedPassword: hashedPassword,
		ID:             userId,
	}

	user, err := cfg.db.UpdateUserEmailAndPassword(req.Context(), updateUserEmailAndPasswordParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
