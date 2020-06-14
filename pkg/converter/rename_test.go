package converter

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestRename_Convert(t *testing.T) {
	type fields struct {
		From string
		To   string
	}
	type args struct {
		json gjson.Result
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    gjson.Result
		wantErr bool
	}{
		{
			"rename existing key",
			fields{From: "foo", To: "bar"},
			args{gjson.Parse("{\"foo\": 1}")},
			gjson.Parse("{\"bar\":1}"),
			false,
		},
		{
			"do nothign if key not exists",
			fields{From: "foo", To: "bar"},
			args{gjson.Parse("{\"foo2\":1}")},
			gjson.Parse("{\"foo2\":1}"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Rename{
				From: tt.fields.From,
				To:   tt.fields.To,
			}
			got, err := c.Convert(tt.args.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenameKey.Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RenameKey.Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
