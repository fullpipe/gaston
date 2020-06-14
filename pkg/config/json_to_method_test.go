package config

import (
	"reflect"
	"testing"

	"github.com/fullpipe/gaston/pkg/converter"
	"github.com/fullpipe/gaston/pkg/remote"
	"github.com/tidwall/gjson"
)

func TestJsonToMethod(t *testing.T) {
	type args struct {
		json gjson.Result
	}
	tests := []struct {
		name    string
		args    args
		want    remote.Method
		wantErr bool
	}{
		{
			name:    "It accepts objects only",
			args:    args{gjson.Parse("[]")},
			want:    remote.Method{},
			wantErr: true,
		},
		{
			name:    "Host is required",
			args:    args{gjson.Parse(`{ "name": "foo" }`)},
			want:    remote.Method{},
			wantErr: true,
		},
		{
			name:    "Name is required",
			args:    args{gjson.Parse(`{ "host": "foo" }`)},
			want:    remote.Method{},
			wantErr: true,
		},
		{
			name:    "It returns basic method",
			args:    args{gjson.Parse(`{ "host": "foo", "name": "bar", "remoteName": "mar" }`)},
			want:    remote.Method{Host: "foo", Name: "bar", RemoteName: "mar"},
			wantErr: false,
		},
		{
			name:    "Remote name could be empty",
			args:    args{gjson.Parse(`{ "host": "foo", "name": "bar" }`)},
			want:    remote.Method{Host: "foo", Name: "bar"},
			wantErr: false,
		},
		{
			name:    "Error on invalid param converter",
			args:    args{gjson.Parse(`{ "host": "foo", "name": "bar", "paramConverters": {} }`)},
			want:    remote.Method{},
			wantErr: true,
		},
		{
			name:    "Error on invalid param converter",
			args:    args{gjson.Parse(`{ "host": "foo", "name": "bar", "paramConverters": [{"from": "foo", "to": "bar"}] }`)},
			want:    remote.Method{},
			wantErr: true,
		},
		{
			name:    "It parse paramConverters",
			args:    args{gjson.Parse(`{ "host": "foo", "name": "bar", "paramConverters": [{"type": "rename", "from": "foo", "to": "bar"}] }`)},
			want:    remote.Method{Host: "foo", Name: "bar", ParamConverters: []converter.Converter{&converter.Rename{From: "foo", To: "bar"}}},
			wantErr: false,
		},
		{
			name:    "Error on invalid result converter",
			args:    args{gjson.Parse(`{ "host": "foo", "name": "bar", "resultConverters": {} }`)},
			want:    remote.Method{},
			wantErr: true,
		},
		{
			name:    "Error on invalid result converter",
			args:    args{gjson.Parse(`{ "host": "foo", "name": "bar", "resultConverters": [{"from": "foo", "to": "bar"}] }`)},
			want:    remote.Method{},
			wantErr: true,
		},
		{
			name:    "It parse resultConverters",
			args:    args{gjson.Parse(`{ "host": "foo", "name": "bar", "resultConverters": [{"type": "rename", "from": "foo", "to": "bar"}] }`)},
			want:    remote.Method{Host: "foo", Name: "bar", ResultConverters: []converter.Converter{&converter.Rename{From: "foo", To: "bar"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JsonToMethod(tt.args.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonToMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonToMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
