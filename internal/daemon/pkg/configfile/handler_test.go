package configfile

import (
	"reflect"
	"testing"
	"time"

	"github.com/hanjunlee/awscred/core"
	"gopkg.in/ini.v1"
)

func LoadFile(raw string) *ini.File {
	f, _ := ini.Load([]byte(raw))
	return f
}

func IniDeepEqual(t *testing.T, f *ini.File, cmp *ini.File) {
	// check sections
	if !reflect.DeepEqual(len(f.SectionStrings()), len(cmp.SectionStrings())) {
		t.Errorf("IniDeepEqual sections is not equal, %v != %v", f.SectionStrings(), cmp.SectionStrings())
	}

	for _, sec := range f.Sections() {
		cs, _ := cmp.GetSection(sec.Name())

		for _, key := range sec.Keys() {
			ck, err := cs.GetKey(key.Name())
			if err != nil {
				t.Errorf("IniDeepEqual key is exist in the cmp: %s", key.Name())
			}

			if !reflect.DeepEqual(key.Value(), ck.Value()) {
				t.Errorf("IniDeepEqual the value of %s is not equal: %s != %s", key.Name(), key.Value(), ck.Value())
			}
		}
	}
}

func TestIniHandler_mapCfgToConfigs(t *testing.T) {
	type fields struct {
		filepath string
	}
	type args struct {
		cfg *ini.File
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]core.Config
	}{
		{
			name: "return empty",
			fields: fields{
				filepath: "",
			},
			args: args{
				cfg: LoadFile(``),
			},
			want: map[string]core.Config{},
		},
		{
			name: "return multiple configurations",
			fields: fields{
				filepath: "",
			},
			args: args{
				cfg: LoadFile(`
[default]
on = true
serial = serial
duration = 3600
aws_access_key_id     = key
aws_secret_access_key = secret
aws_session_token = token
expiration = 2020-08-10T21:00:00Z

[profile]
on = false
serial = serial
duration = 7200
`),
			},
			want: map[string]core.Config{
				"default": {
					On:             true,
					SerialNumber:   "serial",
					DurationSecond: 3600,
					Cache: core.SessionToken{
						AccessKeyID:     "key",
						SecretAccessKey: "secret",
						SessionToken:    "token",
						Expiration:      time.Date(2020, 8, 10, 21, 0, 0, 0, time.UTC),
					},
				},
				"profile": {
					On:             false,
					SerialNumber:   "serial",
					DurationSecond: 7200,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &IniHandler{
				filepath: tt.fields.filepath,
			}
			if got := h.mapCfgToConfigs(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IniHandler.mapCfgToConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIniHandler_mapConfigsToCfg(t *testing.T) {
	type fields struct {
		filepath string
	}
	type args struct {
		configs map[string]core.Config
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
				filepath: "",
			},
			args: args{
				configs: map[string]core.Config{
					"default": {
						On:             true,
						SerialNumber:   "serial",
						DurationSecond: 3600,
						Cache: core.SessionToken{
							AccessKeyID:     "key",
							SecretAccessKey: "secret",
							SessionToken:    "token",
							Expiration:      time.Date(2020, 8, 10, 21, 0, 0, 0, time.UTC),
						},
					},
					"profile": {
						On:             false,
						SerialNumber:   "serial",
						DurationSecond: 7200,
					},
				},
			},
			want: LoadFile(`
[default]
on = true
serial = serial
duration = 3600
aws_access_key_id = key
aws_secret_access_key = secret
aws_session_token = token
expiration = 2020-08-10T21:00:00Z

[profile]
on = false
serial = serial
duration = 7200`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &IniHandler{
				filepath: tt.fields.filepath,
			}
			got := h.mapConfigsToCfg(tt.args.configs)
			IniDeepEqual(t, got, tt.want)
		})
	}
}
