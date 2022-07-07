package controllers

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"social_network_project/entities"
	"testing"
	"time"
)

func TestCreate_CreateToken(t *testing.T) {

	create := Create{}

	id := "6c08496b-b721-4e06-b0b7-1905524c9da2"

	tokenExpected := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := tokenExpected.SignedString([]byte("key"))
	assert.Nil(t, err)

	tokenStructExpected := entities.Token{
		Token: tokenString,
	}

	tokenDecodeExpected := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenStructExpected.Token, tokenDecodeExpected, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})
	assert.Nil(t, err)

	tokenStruct, err := create.Token(id)
	assert.Nil(t, err)

	tokenDecode := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenStruct.Token, tokenDecode, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})

	assert.Equal(t, tokenDecodeExpected["id"], tokenDecode["id"])
}
