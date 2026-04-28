package postgres

import "testing"

func TestVerifyCardToken(t *testing.T) {
	if VerifyCardToken("") {
		t.Fatalf("expected false for empty card token")
	}
	if !VerifyCardToken("tok_1234") {
		t.Fatalf("expected true for non-empty card token")
	}
}
