// Package traefiktokenprocessor provides a Traefik plugin to process JWT tokens
package traefiktokenprocessor

import (
	"context"
	"net/http"
	"strings"
)

// Config holds the plugin configuration.
type Config struct {
	SourceHeader      string `json:"sourceHeader,omitempty"`
	DestinationHeader string `json:"destinationHeader,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		SourceHeader:      "X-Gs-Access-Token",
		DestinationHeader: "X-Api-Apigateway-X-Userinfo",
	}
}

// TokenProcessor is a plugin that extracts a part from a JWT token.
type TokenProcessor struct {
	next   http.Handler
	name   string
	config *Config
}

// New creates a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &TokenProcessor{
		next:   next,
		name:   name,
		config: config,
	}, nil
}

func (t *TokenProcessor) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	token := req.Header.Get(t.config.SourceHeader)

	if token != "" {
		parts := strings.Split(token, ".")
		if len(parts) >= 3 {
			// Get the part between first and second dots (index 1)
			tokenPart := parts[1]

			// Add it as a new header
			req.Header.Set(t.config.DestinationHeader, tokenPart)
		}
	}

	t.next.ServeHTTP(rw, req)
}
