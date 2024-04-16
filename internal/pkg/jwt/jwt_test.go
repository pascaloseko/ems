package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	username := "testuser"
	tokenString, err := GenerateToken(username)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestParseToken(t *testing.T) {
	username := "testuser"
	tokenString, err := GenerateToken(username)
	assert.NoError(t, err)
	username, err = ParseToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, username, "testuser")
}

func TestParseTokenInvalid(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwiaWF0IjoxNjU1OTY1NDI1LCJleHAiOjE2NTU5NzI2MjV9.sdfsdfsdf"
	username, err := ParseToken(tokenString)
	assert.Error(t, err)
	assert.Empty(t, username)
}


func TestErrorUsernameEmptyString(t *testing.T) {
	_, err := GenerateToken("")
	assert.Error(t, err)
}