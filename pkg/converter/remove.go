package converter

import (
	"errors"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Remove struct {
	Name string
}

func (c *Remove) Convert(json gjson.Result) (gjson.Result, error) {
	jsonRaw, _ := sjson.Delete(json.Raw, c.Name)

	return gjson.Parse(jsonRaw), nil
}

func NewRemove(json gjson.Result) (Converter, error) {
	if !json.Get("name").Exists() {
		return nil, errors.New("Specify param name to remove")
	}

	return &Remove{
		Name: json.Get("name").String(),
	}, nil
}
