package database

import (
	"errors"
	"fmt"
	"sort"
)

type Chirp struct {
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
}

func (db *DB) SaveChirp(user_id int, body string) (Chirp, error) {
	data, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(data.Chirps) + 1

	chirp := Chirp{
		Id:       id,
		AuthorId: user_id,
		Body:     body,
	}
	data.Chirps[id] = chirp
	err = db.WriteDB(data)
	if err != nil {
		err_message := fmt.Sprintf("CreateChirp: Couldn't write file: %v", data)
		return Chirp{}, errors.New(err_message)
	}

	return chirp, err
}

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

func (db *DB) GetChirpById(id int) (Chirp, error) {
	data, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}

	for _, chirp := range data.Chirps {
		if chirp.Id == id {
			return chirp, nil
		}
	}
	return Chirp{}, errors.New("no chirp found")
}

func (db *DB) DeleteChirp(id int) error {
	data, err := db.LoadDB()
	if err != nil {
		return err
	}

	delete(data.Chirps, id)
	err = db.WriteDB(data)
	if err != nil {
		return err
	}

	return nil
}
