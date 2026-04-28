package ratelimit

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/time/rate"
)

// ErrRateLimited is the sentinel returned when a client exceeds the rate limit.
// Use errors.Is(err, ErrRateLimited) to detect rate-limit rejection through
// connect.Error wrapping.
var ErrRateLimited = errors.New("rate limited")

type limiterEntry struct{ lim *rate.Limiter }

type ipLimiter struct {
	mu sync.Mutex
	m  map[string]*limiterEntry
	r  rate.Limit
	b  int
}

// NewLoginInterceptor creates a server-side Connect interceptor that rate-limits
// LoginUser calls per-client IP.
func NewLoginInterceptor(rps float64, burst int) connect.Interceptor {
	l := &ipLimiter{m: make(map[string]*limiterEntry), r: rate.Limit(rps), b: burst}
	return &loginLimiter{l: l}
}

type loginLimiter struct{ l *ipLimiter }

func (ll *loginLimiter) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		proc := req.Spec().Procedure
		if !strings.HasSuffix(proc, "/LoginUser") {
			return next(ctx, req)
		}
		ip := clientIP(req)
		if !ll.l.allow(ip) {
			return nil, connect.NewError(connect.CodeResourceExhausted, ErrRateLimited)
		}
		return next(ctx, req)
	}
}

func (ll *loginLimiter) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}
func (ll *loginLimiter) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

func (l *ipLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	ent := l.m[ip]
	if ent == nil {
		ent = &limiterEntry{lim: rate.NewLimiter(l.r, l.b)}
		l.m[ip] = ent
	}
	return ent.lim.AllowN(time.Now(), 1)
}

func clientIP(req connect.AnyRequest) string {
	// Prefer X-Forwarded-For when present
	if xff := req.Header().Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if xr := req.Header().Get("X-Real-IP"); xr != "" {
		return strings.TrimSpace(xr)
	}
	if p := req.Peer(); p.Addr != "" {
		return peerHost(p.Addr)
	}
	return "unknown"
}

// peerHost normalizes a peer address into a stable bucket key for the rate
// limiter. Inputs may be "ip:port", "[ipv6]:port", a bare "ip", or a bare
// "ipv6". SplitHostPort handles the first two; the bare-IP case falls through
// to the raw value, which is already correct (no port to strip).
func peerHost(addr string) string {
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}
