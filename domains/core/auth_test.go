package core

import (
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	signingSecret := "hello world"
	userID, email := uuid.NewString(), gofakeit.Email()
	td := &TokenData{
		UserID: userID,
		Email:  email,
	}

	token, err := td.Generate(signingSecret)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	assert.Len(t, strings.Split(token, "."), 3)
}

func TestDecodeToken(t *testing.T) {
	signingSecret := "hello world"
	userID, email := uuid.NewString(), gofakeit.Email()
	td := &TokenData{
		UserID: userID,
		Email:  email,
	}
	token, err := td.Generate(signingSecret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	tk, err := Decode(token, signingSecret)
	require.NoError(t, err)
	require.NotNil(t, tk)
	assert.Equal(t, td.UserID, tk.UserID)
	assert.Equal(t, td.Email, tk.Email)
}

func TestGeneratePassword(t *testing.T) {
	p := "password"
	h, err := GeneratePassword(p)
	require.NoError(t, err)
	require.NotEmpty(t, h)
}

func TestComparePassword(t *testing.T) {
	p := "password"
	h, err := GeneratePassword(p)
	require.NoError(t, err)
	require.NotEmpty(t, h)

	require.True(t, CheckPasswordHash(p, h))
}
