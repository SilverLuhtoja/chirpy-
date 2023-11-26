package database

import (
	"errors"
	"fmt"
)

func (db *DB) UpdateToMember(id int) error {
	data, err := db.LoadDB()
	if err != nil {
		return err
	}

	user := data.Users[id]

	updateUser := User{
		Id:          user.Id,
		Email:       user.Email,
		Password:    user.Password,
		IsChirpyRed: true,
	}

	data.Users[id] = updateUser
	err = db.WriteDB(data)
	if err != nil {
		err_message := fmt.Sprintf("UpdateToMember: Couldn't write file: %v", data)
		return errors.New(err_message)
	}

	return nil
}
