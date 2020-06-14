package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tokenRaw := req.Header.Get("Authorization")
		tokens := strings.Split(tokenRaw, "Bearer")

		fmt.Println(tokens)
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
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte("qwertyuiopasdfghjklzxcvbnm123456"), nil
		})

		if !token.Valid {
			fmt.Println(err)
			next.ServeHTTP(w, req)
			return
		}

		gastonContext := GetContext(req)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			switch t := claims["roles"].(type) {
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

			gastonContext.Headers["X-Verified-Roles"] = gastonContext.Roles
			gastonContext.Headers["X-Verified-User"] = []string{claims["sub"].(string)}
		}

		req = SetContext(req, gastonContext)

		next.ServeHTTP(w, req)
	})
}
