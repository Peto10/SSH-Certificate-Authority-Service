package server

import (
	"net/http"
	"time"

	"github.com/Peto10/SSH-like-Certificate-Authority-Service/internal/api"
)

const (
	defaultServerReadTimeout       = 1 * time.Second
	defaultServerReadHeaderTimeout = 2 * time.Second
	defaultServerWriteTimeout      = 1 * time.Second
	defaultServerIdleTimeout       = 30 * time.Second
)

func NewServer(c *api.Controller, hostname string) *http.Server {
	router := http.NewServeMux()

	router.HandleFunc("POST /sign", c.Sign)

	return &http.Server{
		Addr:              hostname,
		Handler:           router,
		ReadTimeout:       defaultServerReadTimeout,
		ReadHeaderTimeout: defaultServerReadHeaderTimeout,
		WriteTimeout:      defaultServerWriteTimeout,
		IdleTimeout:       defaultServerIdleTimeout,
	}
}
