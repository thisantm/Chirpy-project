package main

import (
	"io/fs"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

func (db *DB) ensureDB() error {
	_, err := os.OpenFile(db.path, os.O_RDWR|os.O_CREATE, fs.ModePerm)
	if err == os.ErrNotExist {
		log.Print("Creating new Database...")
	} else if err != nil {
		return err
	}

	return nil
}
