package converter

import "github.com/tidwall/gjson"

// Converter can do things with rpc params or results
type Converter interface {
	// Convert does transformations on json
	Convert(json gjson.Result) (gjson.Result, error)
}

// ConverterFactory returns instance of specific Converter.
// Each Converter requires ConverterFactory. ConverterFactory should be
// added to NewProvider function
type ConverterFactory func(json gjson.Result) (Converter, error)
