package converter

import (
	"reflect"
	"testing"

	"github.com/tidwall/gjson"
)

func TestNewSnakeCase(t *testing.T) {
	type args struct {
		json gjson.Result
	}
	tests := []struct {
		name    string
		args    args
		want    Converter
		wantErr bool
	}{
		{
			"it returns error if no name provided",
			args{gjson.Parse(`{"foo":"bar"}`)},
			nil,
			true,
		},
		{
			"it returns replace converter",
			args{gjson.Parse(`{"name":"fooBar"}`)},
			&Rename{From: "fooBar", To: "foo_bar"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSnakeCase(tt.args.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSnakeCase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
