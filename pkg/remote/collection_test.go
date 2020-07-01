package remote

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
	}
	m1 := Method{Name: "foo"}
	m2 := Method{Name: "foo2"}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Method
	}{
		{
			"find method by name",
			fields{Methods: []Method{
				m1, m2,
			}},
			args{methodName: "foo"},
			&m1,
		},
		{
			"returns nil if method not found",
			fields{Methods: []Method{
				m1, m2,
			}},
			args{methodName: "bar"},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MethodCollection{
				Methods: tt.fields.Methods,
			}
			if got := c.Find(tt.args.methodName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MethodCollection.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMethodCollection_Merge(t *testing.T) {
	type fields struct {
		Methods []Method
	}
	type args struct {
		c2 MethodCollection
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   MethodCollection
	}{
		{
			"it merges collections",
			fields{Methods: []Method{Method{Name: "foo"}}},
			args{MethodCollection{Methods: []Method{Method{Name: "bar"}}}},
			MethodCollection{Methods: []Method{Method{Name: "foo"}, Method{Name: "bar"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MethodCollection{
				Methods: tt.fields.Methods,
			}
			if got := c.Merge(tt.args.c2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MethodCollection.Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
