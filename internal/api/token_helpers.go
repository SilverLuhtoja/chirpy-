package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *ApiConfig) createAccessToken(issuer string, id int) (string, error) {
	var new_time_duration time.Duration = time.Hour
	if issuer == REFRESH_ISSUER {
		new_time_duration = time.Hour * time.Duration(24*60) // 60 days
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(new_time_duration)),
		Subject:   fmt.Sprint(id),
	})

	jwt, err := newToken.SignedString([]byte(cfg.JWT))
	if err != nil {
		message := fmt.Sprintf("couldn't sign jwt: %s", err)
		return "", errors.New(message)
	}

	return jwt, nil
}

func (cfg *ApiConfig) createAccessTokenFromHeader(token *jwt.Token) (string, error) {
	user_id, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    ACCESS_ISSUER,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject:   user_id,
	})

	jwt, err := newToken.SignedString([]byte(cfg.JWT))
	if err != nil {
		message := fmt.Sprintf("couldn't sign jwt: %s", err)
		return "", errors.New(message)
	}

	return jwt, nil
}

func (cfg *ApiConfig) getTokenFromHeader(expectedIssure string, r *http.Request) (token *jwt.Token, err error) {
	auth := r.Header.Get("Authorization") //Authorization: Bearer <token>
	if auth == "" {
		return token, errors.New("token is invalid")
	}

	refresh_token := strings.Split(auth, " ")[1]
	token, err = jwt.ParseWithClaims(refresh_token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT), nil
	})
	if err != nil {
		return token, errors.New("not authorized")
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil || issuer != expectedIssure {
		return token, errors.New("problem with issuer")
	}

	return token, nil
}

func (cfg *ApiConfig) getUserIdFromAccessToken(r *http.Request) (zero int, err error) {
	token, err := cfg.getTokenFromHeader(ACCESS_ISSUER, r)
	if err != nil {
		return zero, errors.New("user is not authorized")
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return zero, errors.New("couldn't get issuer from token")
	}

	if issuer != "chirpy-access" {
		return zero, errors.New("wrong token provided")
	}

	user_id, err := token.Claims.GetSubject()
	if err != nil {
		return zero, errors.New("couldn't get id from token")
	}

	intId, err := strconv.Atoi(user_id)
	if err != nil {
		err_message := fmt.Sprintf("getUserIdFromAccessToken: parsing error: %v", user_id)
		return zero, errors.New(err_message)
	}

	return intId, nil
}
