package server

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// analysisLimiter allows 1 analysis per IP per 2 minutes.
var analysisLimiter = &ipRateLimiter{
	ips:     make(map[string]time.Time),
	window:  2 * time.Minute,
}

type ipRateLimiter struct {
	mu     sync.Mutex
	ips    map[string]time.Time
	window time.Duration
}

func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if last, ok := l.ips[ip]; ok && time.Since(last) < l.window {
		return false
	}
	l.ips[ip] = time.Now()
	return true
}

func (l *ipRateLimiter) clear(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.ips, ip)
}

func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.Split(fwd, ",")[0]
	}
	// strip port
	addr := r.RemoteAddr
	if i := strings.LastIndex(addr, ":"); i != -1 {
		return addr[:i]
	}
	return addr
}

func (s *Server) rateLimitAnalysis(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := realIP(r)
		if !analysisLimiter.allow(ip) {
			http.Error(w, "Too many requests — please wait a moment before analyzing another VOD.", http.StatusTooManyRequests)
			return
		}
		next(w, r)
		// if the client cancelled the request, clear the rate limit so they can retry immediately
		if r.Context().Err() != nil {
			analysisLimiter.clear(ip)
		}
	}
}
