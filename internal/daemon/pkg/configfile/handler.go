package configfile

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hanjunlee/awscred/core"
	"gopkg.in/ini.v1"
)

const (
	keyOn                = "on"
	keySerialNumber      = "serial"
	keyDurationSecond    = "duration"
	keyAwsAccessKeyID    = "aws_access_key_id"
	keyAwsSecretAccessID = "aws_secret_access_key"
	keyAwsSessionToken   = "aws_session_token"
	keyExpiration        = "expiration"
)

type (
	// IniHandler handle the yaml-format configuration file.
	IniHandler struct {
		filepath string
	}
)

// NewIniHandler create a new ini handler.
func NewIniHandler(path string) *IniHandler {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatalf("the file doesn't exist: %s", path)
	}

	return &IniHandler{
		filepath: path,
	}
}

// Read read the configuration file.
func (h *IniHandler) Read() (map[string]core.Config, error) {
	cfg, err := ini.Load(h.filepath)
	if err != nil {
		return nil, err
	}

	configs := h.mapCfgToConfigs(cfg)
	return configs, nil
}

func (h *IniHandler) mapCfgToConfigs(cfg *ini.File) map[string]core.Config {
	configs := make(map[string]core.Config)

	for _, sec := range cfg.Sections() {
		// TODO: investigate "DEFAULT" section.
		if sec.Name() == ini.DefaultSection {
			continue
		}

		profile := sec.Name()
		conf := core.Config{}

		for _, key := range sec.Keys() {
			switch key.Name() {
			case keyOn:
				on, err := sec.Key(keyOn).Bool()
				if err != nil {
					on = false
				}

				conf.On = on

			case keySerialNumber:
				conf.SerialNumber = sec.Key(keySerialNumber).String()

			case keyDurationSecond:
				d, err := sec.Key(keyDurationSecond).Int64()
				if err != nil {
					d = -1
				}
				conf.DurationSecond = d

			case keyAwsAccessKeyID:
				conf.Cache.AccessKeyID = sec.Key(keyAwsAccessKeyID).String()

			case keyAwsSecretAccessID:
				conf.Cache.SecretAccessKey = sec.Key(keyAwsSecretAccessID).String()

			case keyAwsSessionToken:
				conf.Cache.SessionToken = sec.Key(keyAwsSessionToken).String()

			case keyExpiration:
				e := sec.Key(keyExpiration).String()
				t, err := time.Parse(time.RFC3339, e)
				if err != nil {
					t = time.Time{}
				}

				conf.Cache.Expiration = t
			}
		}
		configs[profile] = conf
	}
	return configs
}

func (h *IniHandler) Write(c map[string]core.Config) error {
	cfg := h.mapConfigsToCfg(c)
	return cfg.SaveTo(h.filepath)
}

func (h *IniHandler) mapConfigsToCfg(configs map[string]core.Config) *ini.File {
	cfg := ini.Empty()
	
	for profile, conf := range configs {
		sec, _ := cfg.NewSection(profile)

		on := "false"
		if conf.On {
			on = "true"
		} 
		sec.Key(keyOn).SetValue(on)
		sec.Key(keySerialNumber).SetValue(conf.SerialNumber)
		sec.Key(keyDurationSecond).SetValue(strconv.FormatInt(conf.DurationSecond, 10))

		if conf.Cache.AccessKeyID == "" || conf.Cache.SecretAccessKey  == "" {
			continue
		}
		sec.Key(keyAwsAccessKeyID).Comment = "cached token."
		sec.Key(keyAwsAccessKeyID).SetValue(conf.Cache.AccessKeyID)
		sec.Key(keyAwsSecretAccessID).SetValue(conf.Cache.SecretAccessKey)
		sec.Key(keyAwsSessionToken).SetValue(conf.Cache.SessionToken)
		sec.Key(keyExpiration).SetValue(conf.Cache.Expiration.Format(time.RFC3339))
	}
	return cfg
}
