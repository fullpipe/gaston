package converter

import (
	"testing"
)

func TestRenameKey_Convert(t *testing.T) {
	type fields struct {
		From string
		To   string
	}
	type args struct {
		jsonData string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"rename existing key",
			fields{From: "foo", To: "bar"},
			args{jsonData: "{\"foo\": 1}"},
			"{\"bar\":1}",
			false,
		},
		{
			"do nothign if key not exists",
			fields{From: "foo", To: "bar"},
			args{jsonData: "{\"foo2\":1}"},
			"{\"foo2\":1}",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RenameKey{
				From: tt.fields.From,
				To:   tt.fields.To,
			}
			got, err := c.Convert(tt.args.jsonData)
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
