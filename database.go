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
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if !os.IsNotExist(err) {
		return err
	}

	log.Print("Creating new Database...")
	data, err := json.MarshalIndent(DBStructure{}, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateChirp(body string) (chirpValid, error) {

	data, err := os.ReadFile(db.path)
	if err != nil {
		return chirpValid{}, err
	}

	newData := DBStructure{}
	err = json.Unmarshal(data, &newData)
	if err != nil {
		return chirpValid{}, err
	}

	id := 0
	for k := range newData.Chirps {
		id = max(id, k)
	}
	id++

	chirp := chirpValid{
		Id:   id,
		Body: body,
	}

	newData.Chirps[chirp.Id] = chirp

	marshalData, err := json.MarshalIndent(newData, "", "\t")
	if err != nil {
		return chirpValid{}, err
	}

	err = os.WriteFile(db.path, marshalData, os.ModePerm)
	if err != nil {
		return chirpValid{}, err
	}

	return chirp, nil
}
