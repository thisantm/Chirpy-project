package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestCreateChirp(t *testing.T) {
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

func TestGetChirps(t *testing.T) {
	defer os.Remove(testDbPath)

	cases := []struct {
		input    string
		expected struct {
			dbData DBStructure
			chirps []chirpValid
		}
	}{
		{
			input: "Hello, World!",
			expected: struct {
				dbData DBStructure
				chirps []chirpValid
			}{
				dbData: DBStructure{
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
				chirps: []chirpValid{
					{
						Id:   0,
						Body: "Test 0",
					},
					{
						Id:   1,
						Body: "Test 1",
					},
					{
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

		data, err := os.ReadFile(testDbPath)
		if err != nil {
			t.Errorf("Failed to read Database")
		}

		dbData := DBStructure{}
		err = json.Unmarshal(data, &dbData)
		if err != nil {
			t.Errorf("Unsmarshal error %v", err)
		}

		actualChirps, err := db.GetChirps()
		if err != nil {
			t.Errorf("Failed to get Chirps with error: %v", err)
		}

		if !reflect.DeepEqual(actualChirps, cs.expected.chirps) {
			t.Errorf(
				`The expected chirps %v is not equal to the actual chirps %v`,
				cs.expected.chirps,
				actualChirps,
			)
		}

		if !reflect.DeepEqual(dbData, cs.expected.dbData) {
			t.Errorf(
				`The GetChiprs function must not change the database
				the expected db %v is not equal to the actual db %v`,
				cs.expected.dbData,
				dbData.Chirps,
			)
		}
	}
}
