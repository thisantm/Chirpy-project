package main

import (
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestNewDB(t *testing.T) {
	defer os.Remove("database_test.json")
	cases := []struct {
		input    string
		expected *DB
	}{
		{
			input: "database_test.json",
			expected: &DB{
				path: "database_test.json",
				mux:  &sync.RWMutex{},
			},
		},
	}

	for _, cs := range cases {
		err := os.Remove("database_test.json")
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
		if !reflect.DeepEqual(actual, cs.expected) {
			t.Errorf(
				`The expected database %v:
					path: %v
					mux: %v
				is not equal to the database %v:
					path: %v
					mux: %v`,
				cs.expected,
				cs.expected.path,
				cs.expected.mux,
				actual,
				actual.path,
				actual.mux,
			)
		}
	}
}
