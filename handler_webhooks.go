package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/thisantm/Chirpy-project/internal/auth"
)

type PolkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserId string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		log.Printf("failed to get api key: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if strings.Compare(apiKey, cfg.polkaKey) != 0 {
		log.Printf("invalid api key: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(req.Body)
	polkaWebhookRequest := PolkaWebhookRequest{}
	err = decoder.Decode(&polkaWebhookRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if strings.Compare(polkaWebhookRequest.Event, "user.upgraded") != 0 {
		log.Println("event is not valid")
		w.WriteHeader(204)
		return
	}

	userId, err := uuid.Parse(polkaWebhookRequest.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid uuid (failed to parse)", err)
		return
	}

	err = cfg.db.SetUserChirpyRedStatusTrue(req.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
