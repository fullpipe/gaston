package converter

import "github.com/tidwall/gjson"

type Converter interface {
	Convert(json gjson.Result) (gjson.Result, error)
}

type ConverterFactory func(json gjson.Result) (Converter, error)
