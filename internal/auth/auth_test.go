package auth

import (
	"testing"
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
