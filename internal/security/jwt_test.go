package security

import (
    "os"
    "testing"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

func withEnv(k, v string, fn func()) {
    old, ok := os.LookupEnv(k)
    _ = os.Setenv(k, v)
    defer func() {
        if ok { _ = os.Setenv(k, old) } else { _ = os.Unsetenv(k) }
    }()
    fn()
}

func TestVerify_ValidToken(t *testing.T) {
    _ = os.Unsetenv("JWT_ISSUER"); _ = os.Unsetenv("JWT_AUDIENCE")
    withEnv("JWT_SECRET", "s1", func() {
        v, err := NewVerifierFromEnv()
        if err != nil { t.Fatalf("verifier: %v", err) }
        now := time.Now().UTC()
        tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
            IssuedAt: jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(5*time.Minute)),
        })
        s, _ := tok.SignedString([]byte("s1"))
        if _, err := v.Verify(s); err != nil {
            t.Fatalf("verify failed: %v", err)
        }
    })
}

func TestVerify_InvalidSignature(t *testing.T) {
    _ = os.Unsetenv("JWT_ISSUER"); _ = os.Unsetenv("JWT_AUDIENCE")
    withEnv("JWT_SECRET", "s1", func() {
        v, _ := NewVerifierFromEnv()
        now := time.Now().UTC()
        tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
            IssuedAt: jwt.NewNumericDate(now), NotBefore: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(5*time.Minute)),
        })
        s, _ := tok.SignedString([]byte("wrong"))
        if _, err := v.Verify(s); err == nil {
            t.Fatalf("expected error for invalid signature")
        }
    })
}

func TestVerify_IssuerAudience(t *testing.T) {
    _ = os.Unsetenv("JWT_ISSUER"); _ = os.Unsetenv("JWT_AUDIENCE")
    withEnv("JWT_SECRET", "s1", func() {
        _ = os.Setenv("JWT_ISSUER", "iss1")
        _ = os.Setenv("JWT_AUDIENCE", "aud1")
        v, _ := NewVerifierFromEnv()
        now := time.Now().UTC()
        tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
            Issuer: "iss1", Audience: jwt.ClaimStrings{"aud1"},
            IssuedAt: jwt.NewNumericDate(now), NotBefore: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(5*time.Minute)),
        })
        s, _ := tok.SignedString([]byte("s1"))
        if _, err := v.Verify(s); err != nil { t.Fatalf("verify failed: %v", err) }
        // wrong audience
        tok = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "iss1", Audience: jwt.ClaimStrings{"aud2"}, IssuedAt: jwt.NewNumericDate(now), NotBefore: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(5*time.Minute))})
        s, _ = tok.SignedString([]byte("s1"))
        if _, err := v.Verify(s); err == nil { t.Fatalf("expected audience error") }
    })
}

func TestVerify_KeyRingSupportsOldKey(t *testing.T) {
    _ = os.Unsetenv("JWT_ISSUER"); _ = os.Unsetenv("JWT_AUDIENCE")
    // Ring: s2 (active), s1 (old). Verify token signed with old key succeeds.
    withEnv("JWT_SECRETS", "s2,s1", func() {
        v, _ := NewVerifierFromEnv()
        now := time.Now().UTC()
        tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{IssuedAt: jwt.NewNumericDate(now), NotBefore: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(5*time.Minute))})
        s, _ := tok.SignedString([]byte("s1"))
        if _, err := v.Verify(s); err != nil { t.Fatalf("verify with old key failed: %v", err) }
    })
}
