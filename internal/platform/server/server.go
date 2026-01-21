package server

import (
	"log/slog"
	"net/http"
	"os"
)

// ServerConfig contains server configuration.
type Configuration struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

// ListenAndServe creates an http server
func ListenAndServe(address string, handler http.Handler) {
	slog.Info("Listening on http", "address", address)
	err := http.ListenAndServe(address, handler)
	if err != nil {
		slog.Error("ListenAndServe failed", "error", err)
		os.Exit(1)
	}
}
