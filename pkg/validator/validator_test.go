// +build unit

package validator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type character struct {
	Name      string `validate:"required"`
	AnimeFrom string `validate:"required"`
	Email     string `validate:"required,email"`
	Username  string `validate:"required"`
}

func TestDefaultValidator_Struct(t *testing.T) {
	testData := []struct {
		name    string
		char    character
		checkFn func(t *testing.T, err error)
	}{
		{
			name: "Has all required",
			char: character{
				Name:      "Luffy Monkey",
				AnimeFrom: "One Piece",
				Username:  "luffy.monkey",
				Email:     "luffy.monkey@gmail.com",
			},
			checkFn: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Missing Required Field AnimeFrom",
			char: character{
				Name:     "Luffy Monkey",
				Username: "luffy.monkey",
				Email:    "luffy.monkey@gmail.com",
			},
			checkFn: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.IsType(t, ValidationErr{}, err)
			},
		},
		{
			name: "Missing Required Field Username",
			char: character{
				Name:      "Luffy Monkey",
				AnimeFrom: "One Piece",
				Email:     "luffy.monkey@gmail.com",
			},
			checkFn: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.IsType(t, ValidationErr{}, err)
			},
		},
		{
			name: "Missing Required Field Email",
			char: character{
				Name:      "Luffy Monkey",
				AnimeFrom: "One Piece",
				Username:  "luffy.monkey",
			},
			checkFn: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.IsType(t, ValidationErr{}, err)
			},
		},
		{
			name: "Invalid incorrect email format",
			char: character{
				Name:      "Luffy Monkey",
				AnimeFrom: "One Piece",
				Username:  "luffy.monkey",
				Email:     "unknownemailformat",
			},
			checkFn: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.IsType(t, ValidationErr{}, err)
			},
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			err := Struct(td.char)
			td.checkFn(t, err)
		})
	}
}

func TestVar(t *testing.T) {
	testData := []struct {
		name    string
		input   string
		tag     string
		checkFn func(t *testing.T, err error)
	}{
		{
			name:  "Invalid Email",
			input: "luffywithaweirdemail",
			tag:   "required,email",
			checkFn: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.IsType(t, ValidationErr{}, err)
			},
		},
		{
			name:  "Valid Email",
			input: "luffy.monkey@gmail.com",
			tag:   "required,email",
			checkFn: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			err := Var(td.input, td.tag)
			td.checkFn(t, err)
		})
	}
}

func TestStructExcept(t *testing.T) {
	testData := []struct {
		name    string
		char    character
		checkFn func(t *testing.T, err error)
	}{
		{
			name: "all valid",
			char: character{
				Username: "luffy.monkey",
				Email:    "luffy.monkey@gmail.com",
			},
			checkFn: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "missing email",
			char: character{
				Username: "luffy.monkey",
			},
			checkFn: func(t *testing.T, err error) {
				assert.Error(t, err)
				if assert.IsType(t, ValidationErr{}, err) {
					verr := err.(ValidationErr)
					assert.Len(t, verr, 1)
					_, ok := verr["Email"]
					assert.True(t, ok)
				}
			},
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			err := StructExcept(td.char, "Name", "AnimeFrom")
			td.checkFn(t, err)
		})
	}
}
