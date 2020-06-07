package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fullpipe/gaston/converter"
	"github.com/fullpipe/gaston/remote"
)

func main() {
	// TODO: build collection from config file
	collection := remote.MethodCollection{
		Methods: []remote.Method{
			remote.Method{
				Host:   "http://localhost:9091/rpc",
				Name:   "test1",
				Rename: "s1_test",
				Roles:  []string{"asd", "ROLE_USER"},
				ParamConverters: []converter.Converter{
					&converter.RenameKey{
						From: "email_input",
						To:   "email",
					},
				},
			},
			remote.Method{
				Host:   "http://localhost:9091/rpc",
				Name:   "test2",
				Rename: "s1_test2",
				Roles:  []string{"asd", "ROLE_USER"},
				ParamConverters: []converter.Converter{
					&converter.RenameKey{
						From: "email_input",
						To:   "email",
					},
				},
			},
			remote.Method{
				Host:   "http://localhost:9092/rpc",
				Name:   "test3",
				Rename: "s2_test",
				Roles:  []string{"asd", "ROLE_USER"},
				ParamConverters: []converter.Converter{
					&converter.RenameKey{
						From: "email_input",
						To:   "email",
					},
				},
			},
			remote.Method{
				Host:   "http://localhost:9092/rpc",
				Name:   "test4",
				Rename: "s2_test2",
				Roles:  []string{"asd", "ROLE_USER"},
				ParamConverters: []converter.Converter{
					&converter.RenameKey{
						From: "email_input",
						To:   "email",
					},
				},
			},
		},
	}

	// get config from config file?
	tr := &http.Transport{
		MaxIdleConns:    0,
		IdleConnTimeout: 30,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 2,
	}

	server := remote.Server{
		Remote: remote.Remote{
			Methods: collection,
			Client:  client,
		},
	}

	server.Use(LogMiddleware)
	server.Use(AuthenticationMiddleware)
	// TODO: handle errors
	//http.Handle("/", &server)
	//http.Handle("/", AuthenticationMiddleware(LogMiddleware(&server)))
	http.Handle("/", &server)

	// TODO: move port to envars
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req)

		next.ServeHTTP(w, req)
	})
}

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

		gastonContext := remote.GetContext(req)
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
		}

		req = remote.SetContext(req, gastonContext)

		next.ServeHTTP(w, req)
	})
}
