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

func TestValidatePassword(t *testing.T) {
    cases := []struct{
        name string
        pw   string
        ok   bool
    }{
        {"too short", "Ab1!", false},
        {"letters only", "abcdefghi", false},
        {"letters+digits", "abcdefg1", false},
        {"letters+symbols", "abcdefg!", false},
        {"good mix", "Abcd123!", true},
        {"good long", "Very$trongPassw0rd", true},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            err := validatePassword(c.pw)
            if c.ok && err != nil {
                t.Fatalf("expected ok, got %v", err)
            }
            if !c.ok && err == nil {
                t.Fatalf("expected error")
            }
        })
    }
}
