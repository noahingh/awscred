package credfile

import (
	"reflect"
	"testing"

	"github.com/hanjunlee/awscred/core"
	"gopkg.in/ini.v1"
)

func LoadFile(raw string) *ini.File {
	f, _ := ini.Load([]byte(raw))
	return f
}

func TestIniHandler_mapCfgToCreds(t *testing.T) {
	type fields struct {
		readOnly bool
		filepath string
	}
	type args struct {
		cfg *ini.File
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]core.Cred
	}{
		{
			name: "return default profile",
			fields: fields{
				readOnly: true,
				filepath: "",
			},
			args: args{
				cfg: LoadFile(`[default]
aws_access_key_id     = key
aws_secret_access_key = secret`),
			},
			want: map[string]core.Cred{
				"default": {
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
				},
			},
		},
		{
			name: "return multiple profiles",
			fields: fields{
				readOnly: true,
				filepath: "",
			},
			args: args{
				cfg: LoadFile(`
[default]
aws_access_key_id     = key
aws_secret_access_key = secret

[profile]
aws_access_key_id     = key
aws_secret_access_key = secret `),
			},
			want: map[string]core.Cred{
				"default": {
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
				},
				"profile": {
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &IniHandler{
				readOnly: tt.fields.readOnly,
				filepath: tt.fields.filepath,
			}
			if got := h.mapCfgToCreds(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IniHandler.mapCfgToCreds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIniHandler_mapCredsToCfg(t *testing.T) {
	type fields struct {
		readOnly bool
		filepath string
	}
	type args struct {
		creds map[string]core.Cred
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ini.File
	}{
		{
			name: "return default",
			fields: fields{
				readOnly: false,
				filepath: "",
			},
			args: args{
				creds: map[string]core.Cred{
					"default": {
						AccessKeyID:     "key",
						SecretAccessKey: "secret",
					},
				},
			},
			want: LoadFile(`
[default]
aws_access_key_id     = key
aws_secret_access_key = secret`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &IniHandler{
				readOnly: tt.fields.readOnly,
				filepath: tt.fields.filepath,
			}
			got := h.mapCredsToCfg(tt.args.creds)
			// check sections
			for _, sec := range got.Sections() {
				if _, err := tt.want.SectionsByName(sec.Name()); err != nil {
					t.Errorf("IniHandler.mapCredsToCfg() section is not exist: %s", sec.Name())
				}
			}
		})
	}
}
