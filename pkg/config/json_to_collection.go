package config

import (
	"errors"

	"github.com/tidwall/gjson"

	"github.com/fullpipe/gaston/pkg/remote"
)

func JsonToCollection(json gjson.Result) (remote.MethodCollection, error) {
	c := remote.MethodCollection{Methods: []remote.Method{}}

	if !json.IsArray() {
		return remote.MethodCollection{}, errors.New("You should pass an array")
	}

	for _, mj := range json.Array() {
		m, err := JsonToMethod(mj)
		if err != nil {
			return remote.MethodCollection{}, err
		}

		c.Add(m)
	}

	return c, nil
}
