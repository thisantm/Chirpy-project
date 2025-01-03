package main

import (
	"net/http"

	"github.com/thisantm/Chirpy-project/internal/auth"
)

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
