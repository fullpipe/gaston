package converter

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type RenameKey struct {
	From string
	To   string
}

func (c *RenameKey) Convert(jsonData string) (string, error) {
	prevValue := gjson.Get(jsonData, c.From)

	if !prevValue.Exists() {
		return jsonData, nil
	}

	jsonData, _ = sjson.Delete(jsonData, c.From)

	return sjson.SetRaw(jsonData, c.To, prevValue.Raw)
}
