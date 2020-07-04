package converter

import (
	"errors"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type CastNumber struct {
	Name string
}

func (c *CastNumber) Convert(json gjson.Result) (gjson.Result, error) {
	value := json.Get(c.Name)

	if !value.Exists() {
		return json, nil
	}

	if value.Type.String() == "Number" {
		return json, nil
	}

	jsonRaw := json.Raw
	if value.Type.String() == "True" {
		jsonRaw, _ = sjson.Set(jsonRaw, c.Name, 1)
	} else if value.Type.String() == "False" {
		jsonRaw, _ = sjson.Set(jsonRaw, c.Name, 0)
	} else if value.Type.String() == "Null" {
		jsonRaw, _ = sjson.Set(jsonRaw, c.Name, 0)
	} else if value.Type.String() == "String" && strings.TrimSpace(value.String()) == "" {
		jsonRaw, _ = sjson.Set(jsonRaw, c.Name, 0)
	} else if value.Type.String() == "String" {
		i, intErr := strconv.ParseInt(value.String(), 10, 64)
		f, floatErr := strconv.ParseFloat(value.String(), 64)

		if floatErr == nil {
			jsonRaw, _ = sjson.Set(jsonRaw, c.Name, f)
		} else if intErr == nil {
			jsonRaw, _ = sjson.Set(jsonRaw, c.Name, i)
		} else {
			jsonRaw, _ = sjson.Set(jsonRaw, c.Name, nil)
		}
	} else {
		jsonRaw, _ = sjson.Set(jsonRaw, c.Name, nil)
	}

	return gjson.Parse(jsonRaw), nil
}

func NewCastNumber(json gjson.Result) (Converter, error) {
	if !json.Get("name").Exists() {
		return nil, errors.New("Specify param name to cast")
	}

	return &CastNumber{
		Name: json.Get("name").String(),
	}, nil
}
