package config

import (
	"errors"

	"github.com/tidwall/gjson"

	"github.com/fullpipe/gaston/pkg/converter"
	"github.com/fullpipe/gaston/pkg/remote"
)

// todo: move to ConfigParser
var converters *converter.Provider = converter.NewProvider()

func JsonToMethod(json gjson.Result) (remote.Method, error) {
	if !json.IsObject() {
		return remote.Method{}, errors.New("You should pass a json object")
	}

	if !json.Get("host").Exists() {
		return remote.Method{}, errors.New("Host is required")
	}

	if !json.Get("name").Exists() {
		return remote.Method{}, errors.New("Name is required")
	}

	m := remote.Method{
		Host:       json.Get("host").String(),
		Name:       json.Get("name").String(),
		RemoteName: json.Get("remoteName").String(),
	}

	json.Get("roles").ForEach(func(key gjson.Result, role gjson.Result) bool {
		m.Roles = append(m.Roles, role.String())

		return true
	})

	convJson := json.Get("paramConverters")
	if convJson.Exists() {
		if !convJson.IsArray() {
			return remote.Method{}, errors.New("paramConverters should be an array")
		}

		for _, cj := range convJson.Array() {
			conv, err := converters.Get(cj.Get("type").String(), cj)
			if err != nil {
				return remote.Method{}, err
			}

			m.ParamConverters = append(m.ParamConverters, conv)
		}
	}

	resultConvJson := json.Get("resultConverters")
	if resultConvJson.Exists() {
		if !resultConvJson.IsArray() {
			return remote.Method{}, errors.New("resultConverters should be an array")
		}

		for _, rj := range resultConvJson.Array() {
			conv, err := converters.Get(rj.Get("type").String(), rj)
			if err != nil {
				return remote.Method{}, err
			}

			m.ResultConverters = append(m.ResultConverters, conv)
		}
	}

	return m, nil
}
