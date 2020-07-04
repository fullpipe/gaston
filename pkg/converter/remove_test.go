package converter

import (
	"reflect"
	"testing"

	"github.com/tidwall/gjson"
)

func TestRemove_Convert(t *testing.T) {
	type fields struct {
		Name string
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
			"it removes param",
			fields{"foo"},
			args{gjson.Parse(`{"foo":"one","bar":"two"}`)},
			gjson.Parse(`{"bar":"two"}`),
			false,
		},
		{
			"it does noting if param not exists",
			fields{"foo"},
			args{gjson.Parse(`{"bar":"two"}`)},
			gjson.Parse(`{"bar":"two"}`),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Remove{
				Name: tt.fields.Name,
			}
			got, err := c.Convert(tt.args.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("Remove.Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove.Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
