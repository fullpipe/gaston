package server

import (
	"context"
	"net/http"
)

// GastonContextKey is the key for http context
type GastonContextKey struct{}

// GastonContext keeps and collects infomation for Requests
type GastonContext struct {
	Roles   []string
	Headers map[string][]string
}

// GetContext returns context from http.Request
func GetContext(req *http.Request) *GastonContext {
	c := req.Context().Value(&GastonContextKey{})
	if c == nil {
		return &GastonContext{[]string{}, map[string][]string{}}
	}

	return c.(*GastonContext)
}

// SetContext sets context to http.Request
func SetContext(req *http.Request, gastonContext *GastonContext) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), &GastonContextKey{}, gastonContext))
}
