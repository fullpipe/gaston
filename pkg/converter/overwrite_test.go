package converter

import (
	"reflect"
	"testing"

	"github.com/tidwall/gjson"
)

func TestOverwrite_Convert(t *testing.T) {
	type fields struct {
		Name     string
		NewValue interface{}
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
			name:    "It overwites value",
			fields:  fields{"foo", 1},
			args:    args{gjson.Parse(`{"foo":"bar"}`)},
			want:    gjson.Parse(`{"foo":1}`),
			wantErr: false,
		},
		{
			name:    "It does nothing if param not exists",
			fields:  fields{"foo", 1},
			args:    args{gjson.Parse(`{"bar":"bar"}`)},
			want:    gjson.Parse(`{"bar":"bar"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Overwrite{
				Name:     tt.fields.Name,
				NewValue: tt.fields.NewValue,
			}
			got, err := c.Convert(tt.args.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("Overwrite.Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Overwrite.Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOverwrite(t *testing.T) {
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
			name:    "Returns error on no name",
			args:    args{gjson.Parse(`{"foo":"bar"}`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Returns error on empty name",
			args:    args{gjson.Parse(`{"name":"", "newValue":1}`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Returns error on no newValue",
			args:    args{gjson.Parse(`{"name":"foo"}`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Returns overwite converter",
			args:    args{gjson.Parse(`{"name":"foo", "newValue": "bar"}`)},
			want:    &Overwrite{"foo", gjson.Parse(`{"newValue": "bar"}`).Get("newValue").Raw},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOverwrite(tt.args.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOverwrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOverwrite() = %v, want %v", got, tt.want)
			}
		})
	}
}
