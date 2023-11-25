package database

import (
	"errors"
	"fmt"
	"time"
)

type Token struct {
	Id string    `json:"id"`
	S  time.Time `json:"s"`
}

func (db *DB) GetTokens() (tokens []string, err error) {
	data, err := db.LoadDB()
	if err != nil {
		return tokens, err
	}

	for _, token := range data.Tokens {
		tokens = append(tokens, token.Id)
	}

	return tokens, nil
}

func (db *DB) RevokeToken(token string) error {
	data, err := db.LoadDB()
	if err != nil {
		return err
	}

	for _, tk := range data.Tokens {
		if tk.Id == token {
			return errors.New("already revoked")
		}
	}

	id := len(data.Tokens) + 1

	if token == "" {
		return errors.New("no token to revoke")
	}

	data.Tokens[id] = Token{Id: token, S: time.Now()}
	err = db.WriteDB(data)
	if err != nil {
		err_message := fmt.Sprintf("RevokeToken: Couldn't write file: %v", data)
		return errors.New(err_message)
	}

	return nil
}
