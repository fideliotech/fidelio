package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Config struct {
	RequestsPerSecond float64
	BurstSize         int

	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	HandlerTimeout time.Duration

	MaxBodySize int64

	EnableSecurityHeaders bool
}

func DefaultConfig() *Config {
	return &Config{
		RequestsPerSecond: 10.0,
		BurstSize:         20,

		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		HandlerTimeout: 30 * time.Second,

		MaxBodySize: 1 << 20,

		EnableSecurityHeaders: true,
	}
}

type Limiter struct {
	ips        map[string]*rate.Limiter
	lastAccess map[string]time.Time
	mu         *sync.RWMutex
	rate       rate.Limit
	burst      int
	ttl        time.Duration
	ticker     *time.Ticker
}

func NewLimiter(rps float64, burst int) *Limiter {
	limiter := &Limiter{
		ips:        make(map[string]*rate.Limiter),
		lastAccess: make(map[string]time.Time),
		mu:         &sync.RWMutex{},
		rate:       rate.Limit(rps),
		burst:      burst,
		ttl:        1 * time.Hour,
		ticker:     time.NewTicker(1 * time.Hour),
	}

	go limiter.cleanup()

	return limiter
}

func (l *Limiter) cleanup() {
	for range l.ticker.C {
		l.mu.Lock()
		now := time.Now()
		for ip, lastAccess := range l.lastAccess {
			if now.Sub(lastAccess) > l.ttl {
				delete(l.ips, ip)
				delete(l.lastAccess, ip)
			}
		}
		l.mu.Unlock()
	}
}

func (l *Limiter) getLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(l.rate, l.burst)
		l.ips[ip] = limiter
	}

	l.lastAccess[ip] = time.Now()

	return limiter
}

func Stack(handler http.Handler, config *Config) http.Handler {
	if config == nil {
		config = DefaultConfig()
	}

	limiter := NewLimiter(config.RequestsPerSecond, config.BurstSize)

	handler = TimeoutMiddleware(handler, config.HandlerTimeout)
	handler = BodyLimitMiddleware(handler, config.MaxBodySize)
	handler = RateLimitMiddleware(handler, limiter)
	handler = SecurityHeadersMiddleware(handler, config.EnableSecurityHeaders)
	handler = LoggingMiddleware(handler)

	return handler
}

func TimeoutMiddleware(next http.Handler, timeout time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		r = r.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			next.ServeHTTP(w, r)
			close(done)
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
			return
		}
	})
}

func BodyLimitMiddleware(next http.Handler, maxSize int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		}
		next.ServeHTTP(w, r)
	})
}

func RateLimitMiddleware(next http.Handler, limiter *Limiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			ip = forwardedFor
		}

		if !limiter.getLimiter(ip).Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SecurityHeadersMiddleware(next http.Handler, enable bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if enable {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' https://fonts.googleapis.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data:; connect-src 'self'")
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		}
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		logRequest(r, rw.statusCode, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

func logRequest(r *http.Request, statusCode int, duration time.Duration) {
	ip := r.RemoteAddr
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ip = forwardedFor
	}

	logger.Info("%s %s %s %d %s %s",
		ip,
		r.Method,
		r.URL.Path,
		statusCode,
		r.UserAgent(),
		duration,
	)
}
