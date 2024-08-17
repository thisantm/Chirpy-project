package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]chirpValid `json:"chirps"`
}

type chirpValid struct {
	Id   string `json:"id"`
	Body string `json:"body"`
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if !os.IsNotExist(err) {
		return err
	}

	log.Print("Creating new Database...")
	data, _ := json.MarshalIndent(DBStructure{}, "", "\t")
	os.WriteFile(db.path, data, os.ModePerm)

	return nil
}

// func (db *DB) CreateChirp(body string) (chirpValid, error) {
// 	chirp := chirpValid{
// 		Id:   id,
// 		Body: body,
// 	}
// 	data, err := os.ReadFile(db.path)
// 	data = append(data)
// }
