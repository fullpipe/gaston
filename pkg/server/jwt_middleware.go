package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// JWTAuthorizationConfig stores config
type JWTAuthorizationConfig struct {
	Header            string
	Scheme            string
	HmacSecret        string
	RolesClaim        string
	UserClaim         string
	RemoteUserHeader  string
	RemoteRolesHeader string
}

func (c *JWTAuthorizationConfig) normilize() {
	if c.Header == "" {
		c.Header = "Authorization"
	}

	if c.Scheme == "" {
		c.Scheme = "Bearer"
	}

	if c.RolesClaim == "" {
		c.RolesClaim = "roles"
	}

	if c.UserClaim == "" {
		c.UserClaim = "sub"
	}

	if c.RemoteUserHeader == "" {
		c.RemoteUserHeader = "X-Verified-User"
	}

	if c.RemoteRolesHeader == "" {
		c.RemoteRolesHeader = "X-Verified-Roles"
	}
}

// NewJWTAuthorizationMiddleware returns Middleware to handle JWT tokens
func NewJWTAuthorizationMiddleware(config JWTAuthorizationConfig) Middleware {
	config.normilize()

	if config.HmacSecret == "" {
		// todo: move it to jwt.Parse?
		log.Fatalln("Specify HmacSecret for JWTAuthorizationConfig")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			tokenRaw := req.Header.Get(config.Header)
			tokens := strings.Split(tokenRaw, config.Scheme)

			if len(tokens) != 2 {
				next.ServeHTTP(w, req)
				return
			}

			tokenString := strings.TrimSpace(tokens[1])
			if tokenString == "" {
				next.ServeHTTP(w, req)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(config.HmacSecret), nil
			})

			if err != nil {
				next.ServeHTTP(w, req)
				return
			}

			if !token.Valid {
				next.ServeHTTP(w, req)
				return
			}

			gastonContext := GetContext(req)
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				switch t := claims[config.RolesClaim].(type) {
				case []interface{}:
					for _, value := range t {
						gastonContext.Roles = append(gastonContext.Roles, value.(string))
					}
				case []string:
					for _, value := range t {
						gastonContext.Roles = append(gastonContext.Roles, value)
					}
				case string:
					t = strings.TrimSpace(t)
					if t != "" {
						gastonContext.Roles = append(gastonContext.Roles, t)
					}
				}

				gastonContext.Headers[config.RemoteRolesHeader] = gastonContext.Roles
				gastonContext.Headers[config.RemoteUserHeader] = []string{claims[config.UserClaim].(string)}

				// todo: bypass other jwt claims
			}

			req = SetContext(req, gastonContext)

			next.ServeHTTP(w, req)
		})
	}
}
