package gateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// Server is the OpenCrab gateway HTTP/WebSocket server.
type Server struct {
	cfg     *Config
	httpSrv *http.Server
	upgrader websocket.Upgrader
	limiter  *rateLimiter
}

// NewServer creates a new gateway server.
func NewServer(cfg *Config) (*Server, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	limit := cfg.RateLimitPerIP
	if limit <= 0 {
		limit = 60
	}
	limiter := newRateLimiter(float64(limit), limit) // limit req/min per IP

	s := &Server{
		cfg: cfg,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin:     checkOriginLoopback, // Security: only allow loopback origins
		},
		limiter: limiter,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/ws", s.middlewareChain(s.handleWebSocket))

	s.httpSrv = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Address(), cfg.Port),
		Handler:      s.securityHeaders(mux),
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}

	return s, nil
}

// Run starts the server and blocks.
func (s *Server) Run() error {
	log.Info().
		Str("addr", s.httpSrv.Addr).
		Bool("loopback", s.cfg.BindLoopback).
		Msg("gateway listening")

	if s.cfg.TLSCertFile != "" && s.cfg.TLSKeyFile != "" {
		return s.httpSrv.ListenAndServeTLS(s.cfg.TLSCertFile, s.cfg.TLSKeyFile)
	}
	return s.httpSrv.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}

// securityHeaders adds security headers to all responses.
func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}

// middlewareChain applies auth + rate limit.
func (s *Server) middlewareChain(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.limiter.allow(r.RemoteAddr) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		if s.cfg.AuthToken != "" || s.cfg.AuthPassword != "" {
			if !s.authenticate(r) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

func (s *Server) authenticate(r *http.Request) bool {
	if auth := r.Header.Get("Authorization"); auth != "" {
		if len(auth) > 7 && auth[:7] == "Bearer " {
			return constantTimeCompare(auth[7:], s.cfg.AuthToken)
		}
	}
	return false
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OpenCrab Gateway\n"))
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("websocket upgrade failed")
		return
	}
	defer conn.Close()

	conn.SetReadLimit(int64(s.cfg.MaxWebSocketMessageSize))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// TODO: implement protocol handlers (sessions, config, etc.)
	log.Info().Str("remote", r.RemoteAddr).Msg("websocket client connected")
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		_ = message
		// Echo for now; full protocol in follow-up
	}
}
