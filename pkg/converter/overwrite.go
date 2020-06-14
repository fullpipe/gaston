package converter

import (
	"errors"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Overwrite struct {
	Name     string
	NewValue interface{}
}

func (c *Overwrite) Convert(json gjson.Result) (gjson.Result, error) {
	prevValue := json.Get(c.Name)

	if !prevValue.Exists() {
		return json, nil
	}

	jsonRaw, _ := sjson.Set(json.Raw, c.Name, c.NewValue)

	return gjson.Parse(jsonRaw), nil
}

func NewOverwrite(json gjson.Result) (Converter, error) {
	name := json.Get("name")
	if !name.Exists() || name.String() == "" {
		return nil, errors.New("Specify param name to overwite")
	}
	if !json.Get("newValue").Exists() {
		return nil, errors.New("Specify newValue")
	}

	return &Overwrite{
		Name:     name.String(),
		NewValue: json.Get("newValue").Raw,
	}, nil
}
