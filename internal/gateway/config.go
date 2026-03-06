package gateway

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

// Config holds gateway configuration with security-first defaults.
type Config struct {
	// Port is the listening port (default 18789).
	Port int

	// BindLoopback when true binds only to 127.0.0.1 (secure default).
	// When false, binds to 0.0.0.0 (requires explicit operator choice).
	BindLoopback bool

	// AuthToken is the bearer token for API/WebSocket auth (required when exposed).
	AuthToken string

	// AuthPassword enables password-based auth (hashed, never stored plain).
	AuthPassword string

	// TLS cert/key paths for HTTPS (optional; when set, HTTP is disabled).
	TLSCertFile string
	TLSKeyFile  string

	// RateLimitPerIP limits requests per IP per minute (0 = use default 60).
	RateLimitPerIP int

	// MaxWebSocketMessageSize in bytes (default 1MB, prevents DoS).
	MaxWebSocketMessageSize int

	// ReadTimeout, WriteTimeout, IdleTimeout in seconds.
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

// DefaultConfig returns security-hardened defaults.
func DefaultConfig() *Config {
	return &Config{
		Port:                   18789,
		BindLoopback:           true,
		RateLimitPerIP:        60,
		MaxWebSocketMessageSize: 1 << 20, // 1MB
		ReadTimeout:            30,
		WriteTimeout:           30,
		IdleTimeout:            120,
	}
}

// Address returns the bind address string.
func (c *Config) Address() string {
	if c.BindLoopback {
		return "127.0.0.1"
	}
	return "0.0.0.0"
}

// Validate checks config for security issues.
func (c *Config) Validate() error {
	if !c.BindLoopback && c.AuthToken == "" && c.AuthPassword == "" {
		return fmt.Errorf("non-loopback bind requires auth (AuthToken or AuthPassword)")
	}
	if c.MaxWebSocketMessageSize <= 0 || c.MaxWebSocketMessageSize > 16<<20 {
		return fmt.Errorf("MaxWebSocketMessageSize must be 1-16MB")
	}
	return nil
}

// GenerateAuthToken creates a cryptographically secure token.
func GenerateAuthToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// LoadConfigFromEnv loads config from environment variables.
func LoadConfigFromEnv() *Config {
	cfg := DefaultConfig()
	if p := os.Getenv("OPENCRAB_PORT"); p != "" {
		var port int
		if _, err := fmt.Sscanf(p, "%d", &port); err == nil {
			cfg.Port = port
		}
	}
	if os.Getenv("OPENCRAB_BIND_ALL") == "1" || os.Getenv("OPENCRAB_BIND_ALL") == "true" {
		cfg.BindLoopback = false
	}
	if t := os.Getenv("OPENCRAB_AUTH_TOKEN"); t != "" {
		cfg.AuthToken = t
	}
	return cfg
}
