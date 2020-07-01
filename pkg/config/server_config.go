package config

import "github.com/fullpipe/gaston/pkg/server"

type ServerConfig struct {
	Server struct {
		Route string `json:"route"`
		Port  int    `json:"port"`
	} `json:"server"`
	JwtAuthorization server.JWTAuthorizationConfig `json:"jwtAuthorization"`
}

func (s *ServerConfig) Normilize() {
	if s.Server.Route == "" {
		s.Server.Route = "/"
	}

	if s.Server.Port == 0 {
		s.Server.Port = 8080
	}
}
