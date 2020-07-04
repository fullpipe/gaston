package converter

import (
	"reflect"
	"testing"

	"github.com/tidwall/gjson"
)

func TestCastNumber_Convert(t *testing.T) {
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
			"it does nothing if param not exists",
			fields{"foo"},
			args{gjson.Parse(`{"bar": "1"}`)},
			gjson.Parse(`{"bar": "1"}`),
			false,
		},
		{
			"it does nothing if param is number",
			fields{"foo"},
			args{gjson.Parse(`{"foo": 1}`)},
			gjson.Parse(`{"foo": 1}`),
			false,
		},
		{
			"it does nothing if param is number",
			fields{"foo"},
			args{gjson.Parse(`{"foo": 1.1}`)},
			gjson.Parse(`{"foo": 1.1}`),
			false,
		},
		{
			"it converts string to number",
			fields{"foo"},
			args{gjson.Parse(`{"foo": "1"}`)},
			gjson.Parse(`{"foo": 1}`),
			false,
		},
		{
			"string with letters is null",
			fields{"foo"},
			args{gjson.Parse(`{"foo": "a1"}`)},
			gjson.Parse(`{"foo": null}`),
			false,
		},
		{
			"string with letters is null",
			fields{"foo"},
			args{gjson.Parse(`{"foo": "1a"}`)},
			gjson.Parse(`{"foo": null}`),
			false,
		},
		{
			"it parses floats",
			fields{"foo"},
			args{gjson.Parse(`{"foo": "1.1"}`)},
			gjson.Parse(`{"foo": 1.1}`),
			false,
		},
		{
			"float with letters is null",
			fields{"foo"},
			args{gjson.Parse(`{"foo": "1.1a"}`)},
			gjson.Parse(`{"foo": null}`),
			false,
		},
		{
			"it casts true to one",
			fields{"foo"},
			args{gjson.Parse(`{"foo": true}`)},
			gjson.Parse(`{"foo": 1}`),
			false,
		},
		{
			"it casts false to zero",
			fields{"foo"},
			args{gjson.Parse(`{"foo": false}`)},
			gjson.Parse(`{"foo": 0}`),
			false,
		},
		{
			"it casts null to zero",
			fields{"foo"},
			args{gjson.Parse(`{"foo": null}`)},
			gjson.Parse(`{"foo": 0}`),
			false,
		},
		{
			"it casts litiral null to null",
			fields{"foo"},
			args{gjson.Parse(`{"foo": "null"}`)},
			gjson.Parse(`{"foo": null}`),
			false,
		},
		{
			"it casts empty string to zero",
			fields{"foo"},
			args{gjson.Parse(`{"foo": ""}`)},
			gjson.Parse(`{"foo": 0}`),
			false,
		},
		{
			"it casts array to null",
			fields{"foo"},
			args{gjson.Parse(`{"foo": [1]}`)},
			gjson.Parse(`{"foo": null}`),
			false,
		},
		{
			"it casts object to null",
			fields{"foo"},
			args{gjson.Parse(`{"foo": {"a":"b"}}`)},
			gjson.Parse(`{"foo": null}`),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CastNumber{
				Name: tt.fields.Name,
			}
			got, err := c.Convert(tt.args.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("CastNumber.Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CastNumber.Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
