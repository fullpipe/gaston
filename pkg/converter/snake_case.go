package converter

import (
	"errors"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

// @see https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func NewSnakeCase(json gjson.Result) (Converter, error) {
	if !json.Get("name").Exists() {
		return nil, errors.New("Specify param name to convert")
	}

	name := json.Get("name").String()
	snakedName := matchFirstCap.ReplaceAllString(name, "${1}_${2}")
	snakedName = matchAllCap.ReplaceAllString(snakedName, "${1}_${2}")

	return &Rename{
		From: name,
		To:   strings.ToLower(snakedName),
	}, nil
}
