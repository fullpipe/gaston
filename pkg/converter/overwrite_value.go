package converter

import (
	"github.com/tidwall/sjson"
)

type OverwriteValue struct {
	key      string
	newValue interface{}
}

func (c *OverwriteValue) Convert(jsonData string) (string, error) {
	return sjson.Set(jsonData, c.key, c.newValue)
}
