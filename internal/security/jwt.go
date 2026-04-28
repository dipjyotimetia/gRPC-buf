package security

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-buf/internal/config"
)

var (
	ErrMissingSecret   = errors.New("jwt secret missing")
	ErrInvalidToken    = errors.New("invalid token")
	ErrInvalidAudience = errors.New("invalid audience")
	ErrInvalidIssuer   = errors.New("invalid issuer")
)

type Verifier struct {
	SignKey   []byte
	VerifyAll [][]byte
	Issuer    string
	Audience  string
	Leeway    time.Duration
}

func NewVerifierFromEnv() (*Verifier, error) {
	secrets := strings.Split(strings.TrimSpace(os.Getenv("JWT_SECRETS")), ",")
	var keys [][]byte
	for _, s := range secrets {
		s = strings.TrimSpace(s)
		if s != "" {
			keys = append(keys, []byte(s))
		}
	}
	if len(keys) == 0 {
		s := strings.TrimSpace(os.Getenv("JWT_SECRET"))
		if s == "" {
			return nil, ErrMissingSecret
		}
		keys = append(keys, []byte(s))
	}
	v := &Verifier{
		SignKey:   keys[0],
		VerifyAll: keys,
		Issuer:    strings.TrimSpace(os.Getenv("JWT_ISSUER")),
		Audience:  strings.TrimSpace(os.Getenv("JWT_AUDIENCE")),
		Leeway:    30 * time.Second,
	}
	return v, nil
}

// NewVerifierFromConfig constructs a Verifier from SecurityConfig without reading env.
func NewVerifierFromConfig(sec config.SecurityConfig) (*Verifier, error) {
	keys := [][]byte{}
	if s := strings.TrimSpace(sec.JWTSecret); s != "" {
		keys = append(keys, []byte(s))
	}
	if len(keys) == 0 {
		return nil, ErrMissingSecret
	}
	v := &Verifier{
		SignKey:   keys[0],
		VerifyAll: keys,
		Issuer:    strings.TrimSpace(sec.JWTIssuer),
		Audience:  strings.TrimSpace(sec.JWTAudience),
		Leeway:    30 * time.Second,
	}
	return v, nil
}

func (v *Verifier) Verify(tokenString string) (*jwt.RegisteredClaims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}
	var (
		claims    *jwt.RegisteredClaims
		parsed    *jwt.Token
		parseErrs []error
	)
	for _, key := range v.VerifyAll {
		c := &jwt.RegisteredClaims{}
		tok, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
			}
			return key, nil
		})
		if err == nil && tok != nil && tok.Valid {
			claims = c
			parsed = tok
			break
		}
		if err != nil {
			parseErrs = append(parseErrs, err)
		}
	}
	if parsed == nil || !parsed.Valid || claims == nil {
		if len(parseErrs) > 0 {
			return nil, fmt.Errorf("%w: %w", ErrInvalidToken, errors.Join(parseErrs...))
		}
		return nil, ErrInvalidToken
	}
	now := time.Now()
	// Time-based checks with leeway
	if claims.NotBefore != nil && now.Add(v.Leeway).Before(claims.NotBefore.Time) {
		return nil, ErrInvalidToken
	}
	if claims.ExpiresAt != nil && now.After(claims.ExpiresAt.Add(v.Leeway)) {
		return nil, ErrInvalidToken
	}
	// Issuer
	if v.Issuer != "" && claims.Issuer != v.Issuer {
		return nil, ErrInvalidIssuer
	}
	// Audience
	if v.Audience != "" {
		ok := slices.Contains(claims.Audience, v.Audience)
		if !ok {
			return nil, ErrInvalidAudience
		}
	}
	return claims, nil
}

// This package only verifies tokens and returns claims; callers can thread
// claims via their own context if needed.
