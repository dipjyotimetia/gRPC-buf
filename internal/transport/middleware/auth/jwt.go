package auth

import (
    "context"
    "strings"

    "connectrpc.com/connect"
    "github.com/grpc-buf/internal/security"
)

type JWTAuthInterceptor struct {
    v      *security.Verifier
    skip   map[string]bool
    header string
}

// NewJWTAuthInterceptor creates an interceptor that validates Bearer tokens
// for all RPCs except those with procedures listed in skipSuffixes, which are
// matched by HasSuffix (e.g., "/LoginUser").
func NewJWTAuthInterceptor(v *security.Verifier, skipSuffixes []string) connect.Interceptor {
    s := map[string]bool{}
    for _, suf := range skipSuffixes {
        s[suf] = true
    }
    return &JWTAuthInterceptor{v: v, skip: s, header: "Authorization"}
}

func (i *JWTAuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
    return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
        if i.shouldSkip(req.Spec().Procedure) {
            return next(ctx, req)
        }
        token := bearer(req.Header().Get(i.header))
        if token == "" {
            return nil, connect.NewError(connect.CodeUnauthenticated, nil)
        }
        if _, err := i.v.Verify(token); err != nil {
            return nil, connect.NewError(connect.CodeUnauthenticated, err)
        }
        return next(ctx, req)
    }
}

func (i *JWTAuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc { return next }
func (i *JWTAuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc { return next }

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

