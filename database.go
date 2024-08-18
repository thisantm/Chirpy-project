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

func (db *DB) loadDB() (DBStructure, error) {
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbData := DBStructure{}
	err = json.Unmarshal(data, &dbData)
	if err != nil {
		return DBStructure{}, err
	}

	return dbData, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	marshalData, err := json.MarshalIndent(dbStructure, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, marshalData, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateChirp(body string) (chirpValid, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()
	if err != nil {
		return chirpValid{}, err
	}

	id := 0
	for k := range dbData.Chirps {
		id = max(id, k+1)
	}

	chirp := chirpValid{
		Id:   id,
		Body: body,
	}

	if dbData.Chirps == nil {
		dbData.Chirps = make(map[int]chirpValid)
	}

	dbData.Chirps[chirp.Id] = chirp

	err = db.writeDB(dbData)
	if err != nil {
		return chirpValid{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]chirpValid, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbData, err := db.loadDB()
	if err != nil {
		return []chirpValid{}, err
	}

	chirps := []chirpValid{}

	for _, chirp := range dbData.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}
