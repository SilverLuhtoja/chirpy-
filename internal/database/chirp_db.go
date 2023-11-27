package database

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
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

func (db *DB) GetChirps(author_id, sortParam string) ([]Chirp, error) {
	data, err := db.LoadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := []Chirp{}
	for _, val := range data.Chirps {
		chirps = append(chirps, val)
	}

	if author_id != "" {
		id, err := strconv.Atoi(author_id)
		if err != nil {
			return []Chirp{}, err
		}
		chirps = filterChirpsByAuthorId(id, chirps)
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortParam == "desc" {
			return chirps[i].Id > chirps[j].Id
		}
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

func filterChirpsByAuthorId(author_id int, list []Chirp) []Chirp {
	filtered_list := []Chirp{}
	for _, chirp := range list {
		if author_id == chirp.AuthorId {
			filtered_list = append(filtered_list, chirp)
		}
	}

	return filtered_list
}
