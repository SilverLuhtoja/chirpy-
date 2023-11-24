package database

import (
	"github.com/golang-jwt/jwt/v5"
)

func (db *DB) GetTokens() ([]*jwt.Token, error) {
	data, err := db.LoadDB()
	if err != nil {
		return []*jwt.Token{}, err
	}
	return data.Tokens, nil
}
