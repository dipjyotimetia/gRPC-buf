package auth

import (
	"context"
	"strings"

	"connectrpc.com/connect/v2"
	"github.com/grpc-buf/internal/security"
)

type JWTAuthInterceptor struct {
	v      *security.Verifier
	skip   map[string]bool
	header string
}

// NewJWTAuthInterceptor creates a server interceptor that validates Bearer
// tokens for all RPCs except those with procedures listed in skipSuffixes,
// which are matched by HasSuffix (e.g., "/LoginUser"). The check runs on
// request headers only, before any request payload is decoded.
func NewJWTAuthInterceptor(v *security.Verifier, skipSuffixes []string) connect.ServerInterceptor {
	s := map[string]bool{}
	for _, suf := range skipSuffixes {
		s[suf] = true
	}
	i := &JWTAuthInterceptor{v: v, skip: s, header: "Authorization"}
	return i.wrap
}

func (i *JWTAuthInterceptor) wrap(next connect.ServerFunc) connect.ServerFunc {
	return func(ctx context.Context, spec connect.Spec, stream connect.ServerStream) error {
		if i.shouldSkip(spec.Procedure) {
			return next(ctx, spec, stream)
		}
		info, _ := connect.CallInfoForServerContext(ctx)
		token := bearer(info.RequestHeader().Get(i.header))
		if token == "" {
			return connect.NewError(connect.CodeUnauthenticated, "")
		}
		if _, err := i.v.Verify(token); err != nil {
			return connect.NewError(connect.CodeUnauthenticated, err.Error()).WithCause(err)
		}
		return next(ctx, spec, stream)
	}
}

func (i *JWTAuthInterceptor) shouldSkip(proc string) bool {
	for suf := range i.skip {
		if strings.HasSuffix(proc, suf) {
			return true
		}
	}
	return false
}

func bearer(h string) string {
	if h == "" {
		return ""
	}
	const p = "Bearer "
	if strings.HasPrefix(h, p) {
		return strings.TrimSpace(h[len(p):])
	}
	return ""
}
