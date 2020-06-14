package config

import (
	"reflect"
	"testing"

	"github.com/fullpipe/gaston/pkg/remote"
	"github.com/tidwall/gjson"
)

func TestJsonToCollection(t *testing.T) {
	type args struct {
		json gjson.Result
	}
	tests := []struct {
		name    string
		args    args
		want    remote.MethodCollection
		wantErr bool
	}{
		{
			name:    "It returns error if you pass not as array",
			args:    args{json: gjson.Parse("{}")},
			want:    remote.MethodCollection{},
			wantErr: true,
		},
		{
			name:    "It returns methods",
			args:    args{json: gjson.Parse(`[{"host":"foo", "name": "bar"}]`)},
			want:    remote.MethodCollection{Methods: []remote.Method{remote.Method{Host: "foo", Name: "bar"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JsonToCollection(tt.args.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonToCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonToCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}
