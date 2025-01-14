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
