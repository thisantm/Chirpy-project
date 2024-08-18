package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
)

const testDbPath = "database_test.json"

func TestNewDB(t *testing.T) {
	defer os.Remove(testDbPath)

	data, err := json.MarshalIndent(DBStructure{}, "", "\t")
	if err != nil {
		log.Print(err)
	}

	cases := []struct {
		input    string
		expected struct {
			db   *DB
			data []byte
		}
	}{
		{
			input: testDbPath,
			expected: struct {
				db   *DB
				data []byte
			}{
				db: &DB{
					path: testDbPath,
					mux:  &sync.RWMutex{},
				},
				data: data,
			},
		},
	}

	for _, cs := range cases {
		err := os.Remove(testDbPath)
		if err != nil {
			log.Print(err)
		}
		actual, err := NewDB(cs.input)
		if err != nil {
			t.Errorf(
				"A error has occured",
			)
			continue
		}
		if !reflect.DeepEqual(actual, cs.expected.db) {
			t.Errorf(
				`The expected database %v:
				path: %v
				mux: %v
				is not equal to the database %v:
				path: %v
				mux: %v`,
				cs.expected,
				cs.expected.db.path,
				cs.expected.db.mux,
				actual,
				actual.path,
				actual.mux,
			)
		}

		data, err := os.ReadFile(actual.path)
		if err != nil {
			t.Errorf(
				"A error has occured",
			)
		}

		if !reflect.DeepEqual(data, cs.expected.data) {
			t.Errorf(
				`The expected data %v is not equal to the actual data %v`,
				cs.expected.data,
				data,
			)
		}
	}
}

func TestCreateChirp(t *testing.T) {
	defer os.Remove(testDbPath)

	err := os.Remove(testDbPath)
	if err != nil {
		log.Print(err)
	}

	dbTest := DBStructure{
		Chirps: map[int]chirpValid{
			0: {
				Id:   0,
				Body: "Test 0",
			},
			1: {
				Id:   1,
				Body: "Test 1",
			},
		},
	}

	data, err := json.MarshalIndent(dbTest, "", "\t")
	if err != nil {
		log.Print(err)
	}

	err = os.WriteFile(testDbPath, data, os.ModePerm)
	if err != nil {
		t.Errorf("Failed to create Test Database")
	}

	cases := []struct {
		input    string
		expected DBStructure
	}{
		{
			input: "Hello, World!",
			expected: DBStructure{
				Chirps: map[int]chirpValid{
					0: {
						Id:   0,
						Body: "Test 0",
					},
					1: {
						Id:   1,
						Body: "Test 1",
					},
					2: {
						Id:   2,
						Body: "Hello, World!",
					},
				},
			},
		},
	}

	for _, cs := range cases {
		db, err := NewDB(testDbPath)
		if err != nil {
			t.Errorf("Failed to create Test Database")
		}

		actual, err := db.CreateChirp(cs.input)
		if err != nil {
			t.Errorf("Failed to create Chirp with error: %v", err)
		}

		if !reflect.DeepEqual(actual, cs.expected.Chirps[2]) {
			t.Errorf(
				`The expected chirp %v is not equal to the actual chirp %v`,
				cs.expected.Chirps[2],
				actual,
			)
		}

		data, err := os.ReadFile(testDbPath)
		if err != nil {
			t.Errorf("Failed to read Database")
		}

		dbData := DBStructure{}
		err = json.Unmarshal(data, &dbData)
		if err != nil {
			t.Errorf("Unsmarshal error %v", err)
		}

		if !reflect.DeepEqual(dbData, cs.expected) {
			t.Errorf(
				`The expected db %v is not equal to the actual db %v`,
				cs.expected,
				dbData.Chirps,
			)
		}
	}
}
