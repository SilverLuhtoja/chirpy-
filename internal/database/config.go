package database

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

func NewDB(path string) *DB {
	db := &DB{path: path, mux: &sync.RWMutex{}}
	_, err := os.ReadFile(db.path)
	if err != nil {
		db.ensureDB()
	}
	db.removeDataWhenDebug()
	return db
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	f, err := os.Create(db.path)
	if err != nil {
		log.Fatal("couldn't create database")
	}
	defer f.Close()
	return nil
}

// loadDB reads the database file into memory
func (db *DB) LoadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := os.ReadFile(db.path)
	if err != nil {
		err_message := fmt.Sprintf("Couldn't read file: %v", db.path)
		return DBStructure{}, errors.New(err_message)
	}
	structure := DBStructure{}

	if len(data) == 0 {
		return DBStructure{Chirps: make(map[int]Chirp), Users: make(map[int]User)}, nil
	}

	err = json.Unmarshal(data, &structure)
	if err != nil {
		err_message := fmt.Sprintf("Couldn't unmarshal: %v", data)
		return DBStructure{}, errors.New(err_message)
	}

	return structure, nil
}

// writeDB writes the database file to disk
func (db *DB) WriteDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	structure, err := json.Marshal(dbStructure)
	if err != nil {
		err_message := fmt.Sprintf("Couldn't marshal: %v", dbStructure)
		return errors.New(err_message)
	}

	err2 := os.WriteFile(db.path, []byte(structure), os.ModeAppend)
	if err2 != nil {
		err_message := fmt.Sprintf("Couldn't write file: %v", dbStructure)
		return errors.New(err_message)
	}
	return nil
}

func (db *DB) removeDataWhenDebug() error {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		err := db.WriteDB(DBStructure{Chirps: make(map[int]Chirp), Users: make(map[int]User)})
		if err != nil {
			err_message := fmt.Sprintf("Couldn't empty database: %v", err)
			return errors.New(err_message)
		}
	}
	return nil
}
