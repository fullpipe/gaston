package config

import "github.com/fullpipe/gaston/pkg/server"

type ServerConfig struct {
	Server struct {
		Route         string
		Port          int
		MethodsPath   string
		RemoteTimeout int
	}
	Jwt server.JWTAuthorizationConfig
}

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
