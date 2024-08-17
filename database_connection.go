package main

import (
	"sync"
)

type DBStructure struct {
	Chirps map[int]chirpValid `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	db.ensureDB()

	return &db, nil
}
