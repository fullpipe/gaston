package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fullpipe/gaston/pkg/converter"
	"github.com/fullpipe/gaston/pkg/remote"
	"github.com/fullpipe/gaston/pkg/server"
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

	s := server.Server{
		Remote: remote.Remote{
			Methods: collection,
			Client:  client,
		},
	}

	s.Use(LogMiddleware)
	s.Use(server.AuthenticationMiddleware)
	// TODO: handle errors
	//http.Handle("/", &server)
	//http.Handle("/", AuthenticationMiddleware(LogMiddleware(&server)))
	http.Handle("/", &s)

	// TODO: move port to envars
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req)

		next.ServeHTTP(w, req)
	})
}
