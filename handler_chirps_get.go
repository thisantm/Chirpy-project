package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/thisantm/Chirpy-project/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	authorId := req.URL.Query().Get("author_id")
	var userId uuid.NullUUID
	if authorId != "" {
		var err error
		userId.UUID, err = uuid.Parse(authorId)
		userId.Valid = true
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "uuid is not valid", err)
			return
		}
	}

	order := req.URL.Query().Get("sort")
	if order == "" {
		order = "asc"
	}
	dbChirps, err := cfg.db.GetAllChirps(req.Context(),
		database.GetAllChirpsParams{
			UserID: userId,
			Order:  order,
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "uuid is not valid", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
