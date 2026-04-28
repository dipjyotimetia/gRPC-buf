package ratelimit

import (
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestPeerHost(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"ipv4 with port", "192.168.1.1:54321", "192.168.1.1"},
		{"ipv4 bare", "192.168.1.1", "192.168.1.1"},
		{"ipv6 with port", "[2001:db8::1]:54321", "2001:db8::1"},
		{"ipv6 loopback with port", "[::1]:8080", "::1"},
		{"ipv6 bare", "2001:db8::1", "2001:db8::1"},
		{"hostname with port", "client.example.com:9090", "client.example.com"},
		{"hostname bare", "client.example.com", "client.example.com"},
		{"empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := peerHost(c.in)
			if got != c.want {
				t.Fatalf("peerHost(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

// TestPeerHost_BucketStability ensures two requests from the same client on
// different source ports map to the same bucket key. A regression here would
// make the rate limiter useless against a single client whose ephemeral port
// rotates per request.
func TestPeerHost_BucketStability(t *testing.T) {
	a := peerHost("203.0.113.42:51001")
	b := peerHost("203.0.113.42:51002")
	if a != b {
		t.Fatalf("expected stable key across ports; got %q vs %q", a, b)
	}
	v6a := peerHost("[2001:db8::5]:51001")
	v6b := peerHost("[2001:db8::5]:51002")
	if v6a != v6b {
		t.Fatalf("expected stable IPv6 key across ports; got %q vs %q", v6a, v6b)
	}
}

func TestLimiterAllowBurstAndThrottle(t *testing.T) {
	l := &ipLimiter{m: make(map[string]*limiterEntry), r: rate.Limit(2), b: 2}
	ip := "1.2.3.4"
	// Burst should allow first two instantly. Evaluate both calls up front so
	// the second isn't short-circuited away when the first fails.
	first := l.allow(ip)
	second := l.allow(ip)
	if !first || !second {
		t.Fatalf("expected initial burst to be allowed; got first=%v second=%v", first, second)
	}
	// Third should sometimes be denied immediately depending on rps; to make deterministic, wait a bit
	if l.allow(ip) {
		// If allowed, immediately call again which should likely be denied
		if l.allow(ip) {
			t.Fatalf("expected throttle after burst")
		}
	} else {
		// denied as expected; after some time, should allow again
		time.Sleep(600 * time.Millisecond)
		if !l.allow(ip) {
			t.Fatalf("expected allow after refill")
		}
	}
}
