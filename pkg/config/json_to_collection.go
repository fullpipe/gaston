package config

import (
	"errors"

	"github.com/tidwall/gjson"

	"github.com/fullpipe/gaston/pkg/remote"
)

// JsonToCollection converts gjson collection to collection of Methods
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

// JsonStringToCollection converts json string to collection of Methods
func JsonStringToCollection(json string) (remote.MethodCollection, error) {
	return JsonToCollection(gjson.Parse(json))
}
