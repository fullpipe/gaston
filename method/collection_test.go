package method

import (
	"reflect"
	"testing"
)

func TestMethodCollection_Find(t *testing.T) {
	type fields struct {
		Methods []Method
	}
	type args struct {
		methodName string
		version    string
	}
	m1 := Method{Name: "foo", Version: "v1"}
	m2 := Method{Name: "foo", Version: "v2"}
	m3 := Method{Name: "foo2"}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Method
	}{
		{
			"find method by name and version",
			fields{Methods: []Method{
				m1, m2,
			}},
			args{methodName: "foo", version: "v2"},
			&m2,
		},
		{
			"find method by both name and version",
			fields{Methods: []Method{
				m1, m2,
			}},
			args{methodName: "foo", version: ""},
			nil,
		},
		{
			"returns nil if method not found",
			fields{Methods: []Method{
				m1, m2,
			}},
			args{methodName: "bar", version: "v2"},
			nil,
		},
		{
			"find versionless method",
			fields{Methods: []Method{
				m1, m2, m3,
			}},
			args{methodName: "foo2"},
			&m3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MethodCollection{
				Methods: tt.fields.Methods,
			}
			if got := c.Find(tt.args.methodName, tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MethodCollection.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
