package main

import (
	"sync"
)

func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	db.ensureDB()

	return &db, nil
}
