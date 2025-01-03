package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/thisantm/Chirpy-project/internal/auth"
)

type RefreshResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	fmt.Println(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "authorization header missing", err)
		return
	}

	dbUserId, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not found", err)
	}

	accessToken, err := auth.MakeJWT(dbUserId, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create access token", err)
	}

	respondWithJSON(w, http.StatusOK, RefreshResponse{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "authorization header missing", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to revoke", err)
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
