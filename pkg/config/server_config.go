package config

import "github.com/fullpipe/gaston/pkg/server"

// ServerConfig holds server configuration
type ServerConfig struct {
	Server struct {
		Route         string
		Port          int
		MethodsPath   string
		RemoteTimeout int
	}
	Jwt server.JWTAuthorizationConfig
}

// Normilize sets up default values for ServerConfig
func (s *ServerConfig) Normilize() {
	if s.Server.Route == "" {
		s.Server.Route = "/"
	}

	if s.Server.Port == 0 {
		s.Server.Port = 8080
	}
	if s.Server.RemoteTimeout == 0 {
		s.Server.RemoteTimeout = 5
	}
}
