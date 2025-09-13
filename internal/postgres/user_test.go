package postgres

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	pw := "SuperSecret123!"
	hashed, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if hashed == pw {
		t.Fatalf("hash should not equal original password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw)); err != nil {
		t.Fatalf("bcrypt compare failed: %v", err)
	}
}
