package converter

import "testing"

func TestOverwriteValue_Convert(t *testing.T) {
	type fields struct {
		key      string
		newValue interface{}
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
			"overwrite value for key",
			fields{key: "foo", newValue: "bar"},
			args{jsonData: "{\"foo\": 1}"},
			"{\"foo\": \"bar\"}",
			false,
		},
		{
			"add key with value if not exists",
			fields{key: "foo", newValue: "bar"},
			args{jsonData: "{\"foo2\":1}"},
			"{\"foo2\":1,\"foo\":\"bar\"}",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &OverwriteValue{
				key:      tt.fields.key,
				newValue: tt.fields.newValue,
			}
			got, err := c.Convert(tt.args.jsonData)
			if (err != nil) != tt.wantErr {
				t.Errorf("OverwriteValue.Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OverwriteValue.Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
