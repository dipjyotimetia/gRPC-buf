package postgres

import "testing"

func TestVerifyCard(t *testing.T) {
	if VerifyCard(0) {
		t.Fatalf("expected false for 0 card number")
	}
	if !VerifyCard(1234) {
		t.Fatalf("expected true for non-zero card number")
	}
}
