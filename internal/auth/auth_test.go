package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashing(t *testing.T) {
	password := "asd"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error on HashPassword %v", err)
		return
	}

	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("Hash not match")
		return
	} else {
		t.Log("Hash Ok!")
	}
}

func TestAuth(t *testing.T) {
	type testCase struct {
		userId              uuid.UUID
		tokenSecret         string
		tokenSecretValidate string
		expiresIn           time.Duration
		wait                time.Duration
		expected            bool
		description         string
	}

	cases := []testCase{
		{
			userId:              uuid.New(),
			tokenSecret:         "secret",
			tokenSecretValidate: "secret",
			expiresIn:           1 * time.Second,
			wait:                0 * time.Second,
			description:         "accepted case",
			expected:            true,
		},
		{
			userId:              uuid.New(),
			tokenSecret:         "secret",
			tokenSecretValidate: "secret",
			expiresIn:           10 * time.Millisecond,
			wait:                20 * time.Millisecond,
			description:         "timeout case",
			expected:            false,
		},
		{
			userId:              uuid.New(),
			tokenSecret:         "secret",
			tokenSecretValidate: "notsecret",
			expiresIn:           1 * time.Second,
			wait:                0 * time.Second,
			description:         "wrong key case",
			expected:            false,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			token, err := MakeJWT(c.userId, c.tokenSecret, c.expiresIn)
			if err != nil {
				t.Fatalf("Error on MakeJWT %v", err)
			}

			time.Sleep(c.wait)

			validatedID, err := ValidateJWT(token, c.tokenSecretValidate)

			if err != nil {
				s := fmt.Sprintf("test case: %s > JWT rejected: %v", c.description, err)
				logResult(t, err, c.expected, s)
			} else if c.userId != validatedID {
				s := fmt.Sprintf("test case: %s > UUID don't match %v\nexpected: %v\ngot: %v", c.description, err, c.userId, validatedID)
				logResult(t, err, c.expected, s)
			} else {
				s := fmt.Sprintf("test case: %s > JWT accepted", c.description)
				logResult(t, err, c.expected, s)
			}
		})
	}
}

func logResult(t *testing.T, err error, expected bool, msg string) {
	if (err == nil && expected == true) || (err != nil && expected == false) {
		t.Logf("Success! %v", msg)
	} else {
		t.Errorf("Failed! %v", msg)
	}
}
