package converter

import (
	"errors"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Rename struct {
	From string
	To   string
}

func (c *Rename) Convert(json gjson.Result) (gjson.Result, error) {
	prevValue := json.Get(c.From)

	if !prevValue.Exists() {
		return json, nil
	}

	jsonRaw, _ := sjson.Delete(json.Raw, c.From)
	jsonRaw, _ = sjson.SetRaw(jsonRaw, c.To, prevValue.Raw)

	return gjson.Parse(jsonRaw), nil
}

func NewRename(json gjson.Result) (Converter, error) {
	if !json.Get("from").Exists() || !json.Get("to").Exists() {
		return nil, errors.New("Specify from and to ")
	}

	return &Rename{
		From: json.Get("from").String(),
		To:   json.Get("to").String(),
	}, nil
}
