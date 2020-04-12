package method

import (
	"testing"

	"github.com/fullpipe/gaston/converter"
)

func TestMethod_IsGranted(t *testing.T) {
	type fields struct {
		Host             string
		Version          string
		Name             string
		Rename           string
		Roles            []string
		ParamConverters  []converter.Converter
		ResultConverters []converter.Converter
	}
	type args struct {
		roles []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"granted if has one of required roles",
			fields{Roles: []string{"ROLE_USER", "ROLE_ADMIN"}},
			args{roles: []string{"ROLE_USER"}},
			true,
		},
		{
			"is not granted if no roles from required",
			fields{Roles: []string{"ROLE_ADMIN"}},
			args{roles: []string{"ROLE_USER"}},
			false,
		},
		{
			"granted if no roles required",
			fields{},
			args{roles: []string{"ROLE_USER"}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Method{
				Host:             tt.fields.Host,
				Version:          tt.fields.Version,
				Name:             tt.fields.Name,
				Rename:           tt.fields.Rename,
				Roles:            tt.fields.Roles,
				ParamConverters:  tt.fields.ParamConverters,
				ResultConverters: tt.fields.ResultConverters,
			}
			if got := m.IsGranted(tt.args.roles); got != tt.want {
				t.Errorf("Method.IsGranted() = %v, want %v", got, tt.want)
			}
		})
	}
}
