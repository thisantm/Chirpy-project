package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
)

var startState = []DBStructure{
	{
		Chirps: map[int]chirpValid{
			1: {
				Id:   1,
				Body: "Test 0",
			},
			2: {
				Id:   2,
				Body: "Test 1",
			},
		},
	},
	{},
}

func TestCreateChirp(t *testing.T) {
	defer os.Remove(testDbPath)

	cases := []struct {
		input    string
		expected struct {
			db    DBStructure
			chirp chirpValid
		}
	}{
		{
			input: "Hello, World!",
			expected: struct {
				db    DBStructure
				chirp chirpValid
			}{
				db: DBStructure{
					Chirps: map[int]chirpValid{
						1: {
							Id:   1,
							Body: "Test 0",
						},
						2: {
							Id:   2,
							Body: "Test 1",
						},
						3: {
							Id:   3,
							Body: "Hello, World!",
						},
					},
				},
				chirp: chirpValid{
					Id:   3,
					Body: "Hello, World!",
				},
			},
		},
		{
			input: "Hello, World!",
			expected: struct {
				db    DBStructure
				chirp chirpValid
			}{
				db: DBStructure{
					Chirps: map[int]chirpValid{
						1: {
							Id:   1,
							Body: "Hello, World!",
						},
					},
				},
				chirp: chirpValid{
					Id:   1,
					Body: "Hello, World!",
				},
			},
		},
	}

	if len(cases) != len(startState) {
		t.Errorf(`The number of test cases: %v
		must be equal to the number of start states: %v`, len(cases), len(startState))
	}

	for i, cs := range cases {
		err := os.Remove(testDbPath)
		if err != nil {
			log.Print(err)
		}

		dbTest := startState[i]

		data, err := json.MarshalIndent(dbTest, "", "\t")
		if err != nil {
			log.Print(err)
		}

		err = os.WriteFile(testDbPath, data, os.ModePerm)
		if err != nil {
			t.Errorf("Failed to create Test Database")
		}

		db, err := NewDB(testDbPath)
		if err != nil {
			t.Errorf("Failed to connect to Test Database")
		}

		actual, err := db.CreateChirp(cs.input)
		if err != nil {
			t.Errorf("Failed to create Chirp with error: %v", err)
		}

		if !reflect.DeepEqual(actual, cs.expected.chirp) {
			t.Errorf(
				`The expected chirp %v is not equal to the actual chirp %v`,
				cs.expected.chirp,
				actual,
			)
		}

		data, err = os.ReadFile(testDbPath)
		if err != nil {
			t.Errorf("Failed to read Database")
		}

		dbData := DBStructure{}
		err = json.Unmarshal(data, &dbData)
		if err != nil {
			t.Errorf("Unsmarshal error %v", err)
		}

		if !reflect.DeepEqual(dbData, cs.expected.db) {
			t.Errorf(
				`The expected db %v is not equal to the actual db %v`,
				cs.expected,
				dbData,
			)
		}
	}
}

func TestGetChirps(t *testing.T) {
	defer os.Remove(testDbPath)

	cases := []struct {
		expected struct {
			dbData DBStructure
			chirps []chirpValid
		}
	}{
		{
			expected: struct {
				dbData DBStructure
				chirps []chirpValid
			}{
				dbData: DBStructure{
					startState[0].Chirps,
				},
				chirps: []chirpValid{
					startState[0].Chirps[1],
					startState[0].Chirps[2],
				},
			},
		},
		{
			expected: struct {
				dbData DBStructure
				chirps []chirpValid
			}{
				dbData: DBStructure{},
				chirps: []chirpValid{},
			},
		},
	}

	if len(cases) != len(startState) {
		t.Errorf("The number of test cases must be equal to the number of start states")
	}

	for i, cs := range cases {
		err := os.Remove(testDbPath)
		if err != nil {
			log.Print(err)
		}

		dbTest := startState[i]

		data, err := json.MarshalIndent(dbTest, "", "\t")
		if err != nil {
			log.Print(err)
		}

		err = os.WriteFile(testDbPath, data, os.ModePerm)
		if err != nil {
			t.Errorf("Failed to create Test Database")
		}

		db, err := NewDB(testDbPath)
		if err != nil {
			t.Errorf("Failed to connect to Test Database")
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

func TestGetChirpById(t *testing.T) {
	defer os.Remove(testDbPath)

	cases := []struct {
		input    int
		expected struct {
			dbData DBStructure
			chirp  chirpValid
		}
	}{
		{
			input: 1,
			expected: struct {
				dbData DBStructure
				chirp  chirpValid
			}{
				dbData: DBStructure{
					startState[0].Chirps,
				},
				chirp: chirpValid{
					Id:   1,
					Body: "Test 0",
				},
			},
		},
		{
			input: 1,
			expected: struct {
				dbData DBStructure
				chirp  chirpValid
			}{
				dbData: DBStructure{},
				chirp:  chirpValid{},
			},
		},
	}

	if len(cases) != len(startState) {
		t.Errorf("The number of test cases must be equal to the number of start states")
	}

	for i, cs := range cases {
		err := os.Remove(testDbPath)
		if err != nil {
			log.Print(err)
		}

		dbTest := startState[i]

		data, err := json.MarshalIndent(dbTest, "", "\t")
		if err != nil {
			log.Print(err)
		}

		err = os.WriteFile(testDbPath, data, os.ModePerm)
		if err != nil {
			t.Errorf("Failed to create Test Database")
		}

		db, err := NewDB(testDbPath)
		if err != nil {
			t.Errorf("Failed to connect to Test Database")
		}

		dbData := DBStructure{}
		err = json.Unmarshal(data, &dbData)
		if err != nil {
			t.Errorf("Unsmarshal error %v", err)
		}

		actualChirp, err := db.GetChirpById(cs.input)
		if err != nil && i == 0 {
			t.Errorf("Failed to get Chirp with error: %v", err)
		}
		if err != nil && i == 1 {
			if err.Error() != "not found" {
				t.Errorf("Failed to get Chirp with error: %v", err)
			}
		}

		if !reflect.DeepEqual(actualChirp, cs.expected.chirp) {
			t.Errorf(
				`The expected chirp %v is not equal to the actual chirp %v`,
				cs.expected.chirp,
				actualChirp,
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
