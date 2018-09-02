package helper

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	testPasswords := []string{
		"Blabla2345",
		"&)*ç(%%çdösfkSDfkj345",
		"SDDfkjwleg*ç34234!$",
	}

	for _, testPassword := range testPasswords {
		hash, err := HashPassword(testPassword)
		if err != nil {
			t.Error(err)
		}

		if !CheckPasswordHash(testPassword, hash) {
			t.Errorf("Password: %s isn't valid!", testPassword)
		}
	}
}

func TestCheckPasswordHash(t *testing.T) {
	type TestMatrix struct {
		ClearTextPassword string
		Hash string
		Correct bool
	}

	testMatrix := []TestMatrix{
		{
			ClearTextPassword: ")sdDfjl=+BBb091!",
			Hash: "$2a$14$F.jfRxIxMysxJA2nQv4zhuwq97hfdNBoKsRis0wy1edesof48o6sO",
			Correct: true,
		},
		{
			ClearTextPassword: "oaskfRWE%çç)df093!",
			Hash: "$2a$14$lELHxPz4dMPcMfFg7HVY8OyWjIuLBz/2dDeoQLL6CRiymcxpT/8um",
			Correct: true,
		},
		{
			ClearTextPassword: "!ldsfSDGJwkelg$)*",
			Hash: "$2a$14$0pxZWwT6Y2lvK4iMZqduEeWGDOngsRwVDVoXGKfliz6YTSmhY1bpi",
			Correct: false,
		},
	}

	for _, test := range testMatrix {
		result := CheckPasswordHash(test.ClearTextPassword, test.Hash)

		if result != test.Correct {
			t.Errorf("Password check on %s was not expected result", test.ClearTextPassword)
		}
	}
}