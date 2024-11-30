package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	cfg.fileServerHits.Store(0)
	cfg.db.Reset(req.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Number of hits to /app* reset and database reset to initial state successfully"))
}
