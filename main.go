package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/fullpipe/gaston/pkg/config"
	"github.com/fullpipe/gaston/pkg/converter"
	"github.com/fullpipe/gaston/pkg/remote"
	"github.com/fullpipe/gaston/pkg/server"
)

func main() {
	files, err := filepath.Glob("methods/*.json")
	if err != nil {
		log.Fatalln(err)
	}

	collection := remote.NewMethodCollection()
	for _, f := range files {
		temp, _ := ioutil.ReadFile(f)
		c, err := config.JsonStringToCollection(string(temp))
		if err != nil {
			log.Fatalln(err)
		}

		collection = collection.Merge(c)
	}

	fmt.Println(collection)

	// TODO: build collection from config file
	c2 := remote.MethodCollection{
		Methods: []remote.Method{
			remote.Method{
				Host:       "http://localhost:9092/rpc",
				Name:       "test5",
				RemoteName: "s2_test2",
				Roles:      []string{"asd", "ROLE_USER"},
				ParamConverters: []converter.Converter{
					&converter.Rename{
						From: "email_input",
						To:   "email",
					},
				},
			},
		},
	}

	collection = collection.Merge(c2)

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
	s.Use(server.NewJWTAuthorizationMiddleware(server.JWTAuthorizationConfig{
		HmacSecret: "qwertyuiopasdfghjklzxcvbnm123456",
	}))
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
