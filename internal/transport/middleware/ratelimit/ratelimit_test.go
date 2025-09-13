package ratelimit

import (
    "testing"
    "time"

    "golang.org/x/time/rate"
)

func TestLimiterAllowBurstAndThrottle(t *testing.T) {
    l := &ipLimiter{m: make(map[string]*limiterEntry), r: rate.Limit(2), b: 2}
    ip := "1.2.3.4"
    // Burst should allow first two instantly
    if !l.allow(ip) || !l.allow(ip) {
        t.Fatalf("expected initial burst to be allowed")
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

