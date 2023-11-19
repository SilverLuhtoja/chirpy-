package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{path: path, mux: &sync.RWMutex{}}
	_, err := os.ReadFile(db.path)
	if err != nil {
		db.ensureDB()
		return db, err
	}
	return db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	f, err := os.Create(db.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(data.Chirps) + 1

	chirp := Chirp{
		Id:   id,
		Body: body,
	}
	data.Chirps[id] = chirp
	err = db.writeDB(data)
	if err != nil {
		err_message := fmt.Sprintf("CreateChirp: Couldn't write file: %v", data)
		return Chirp{}, errors.New(err_message)
	}

	return chirp, err
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := []Chirp{}
	for _, val := range data.Chirps {
		chirps = append(chirps, val)
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})
	return chirps, nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := os.ReadFile(db.path)
	if err != nil {
		err_message := fmt.Sprintf("Couldn't read file: %v", db.path)
		return DBStructure{}, errors.New(err_message)
	}
	structure := DBStructure{}

	if len(data) == 0 {
		return DBStructure{Chirps: make(map[int]Chirp)}, nil
	}

	err = json.Unmarshal(data, &structure)
	if err != nil {
		err_message := fmt.Sprintf("Couldn't unmarshal: %v", data)
		return DBStructure{}, errors.New(err_message)
	}

	return structure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
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
