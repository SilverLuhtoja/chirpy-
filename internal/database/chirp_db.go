package database

import (
	"errors"
	"fmt"
	"sort"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) SaveChirp(body string) (Chirp, error) {
	data, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(data.Chirps) + 1

	chirp := Chirp{
		Id:   id,
		Body: body,
	}
	data.Chirps[id] = chirp
	err = db.WriteDB(data)
	if err != nil {
		err_message := fmt.Sprintf("CreateChirp: Couldn't write file: %v", data)
		return Chirp{}, errors.New(err_message)
	}

	return chirp, err
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	data, err := db.LoadDB()
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
