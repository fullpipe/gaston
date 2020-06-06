package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fullpipe/gaston/converter"
	"github.com/fullpipe/gaston/remote"
)

func main() {
	// TODO: build collection from config file
	collection := remote.MethodCollection{
		Methods: []remote.Method{
			remote.Method{
				Host:   "http://localhost:9090/rpc",
				Name:   "test1",
				Rename: "cohort.getIds",
				Roles:  []string{"asd", "ROLE_USER"},
				ParamConverters: []converter.Converter{
					&converter.RenameKey{
						From: "slug2",
						To:   "slug",
					},
				},
			},
			remote.Method{
				Host:   "http://localhost:9090/rpc",
				Name:   "test2",
				Rename: "cohort.getIds",
				Roles:  []string{"asd", "ROLE_USER"},
				ParamConverters: []converter.Converter{
					&converter.RenameKey{
						From: "slug2",
						To:   "slug",
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
	// TODO: handle errors
	http.HandleFunc("/", server.HandleRequest)

	// TODO: move port to envars
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
