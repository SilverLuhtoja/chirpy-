package database

import (
	"errors"
	"fmt"
	"sort"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Password []byte `json:"password,omitempty"`
	Email    string `json:"email"`
}

func (db *DB) CreateUser(password, email string) (User, error) {
	data, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}
	id := len(data.Users) + 1

	encrypted_pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err_message := fmt.Sprintf("CreateUser: Couldn't encrypt password: %v", data)
		return User{}, errors.New(err_message)
	}

	user := User{
		Id:       id,
		Password: encrypted_pass,
		Email:    email,
	}

	data.Users[id] = user
	err = db.WriteDB(data)
	if err != nil {
		err_message := fmt.Sprintf("CreateUser: Couldn't write file: %v", data)
		return User{}, errors.New(err_message)
	}

	user.Password = []byte("")
	return user, err
}

func (db *DB) UpdateUser(id int, password, email string) (User, error) {
	data, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	encrypted_pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err_message := fmt.Sprintf("UpdateUser: Couldn't encrypt password: %v", data)
		return User{}, errors.New(err_message)
	}

	user := User{
		Id:       id,
		Email:    email,
		Password: encrypted_pass,
	}

	data.Users[id] = user
	err = db.WriteDB(data)
	if err != nil {
		err_message := fmt.Sprintf("CreateUser: Couldn't write file: %v", data)
		return User{}, errors.New(err_message)
	}

	user.Password = []byte("")
	return user, err
}

// GetUsers returns all chirps in the database
func (db *DB) GetUsers() ([]User, error) {
	data, err := db.LoadDB()
	if err != nil {
		return []User{}, err
	}

	users := []User{}
	for _, val := range data.Users {
		users = append(users, val)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})

	return users, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	data, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range data.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("not found")
}
