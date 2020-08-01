package credfile

import (
	"fmt"
	"log"
	"os"

	"github.com/hanjunlee/awscred/core"
	"gopkg.in/ini.v1"
)

const (
	keyAwsAccessKeyID    = "aws_access_key_id"
	keyAwsSecretAccessID = "aws_secret_access_key"
	keyAwsSessionToken   = "aws_session_token"
)

var (
	// ErrReadOnly is the read-only error.
	ErrReadOnly = fmt.Errorf("this handler is read only")
)

// IniHandler handles the credential file.
type IniHandler struct {
	readOnly bool
	filepath string
}

// NewIniHandler create a new credential handler.
func NewIniHandler(readOnly bool, path string) *IniHandler {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatalf("the file doesn't exist: %s", path)
	}

	return nil
}

// Read read the credential file.
func (h *IniHandler) Read() (map[string]core.Cred, error) {
	cfg, err := ini.Load(h.filepath)
	if err != nil {
		return nil, err
	}

	creds := h.mapCfgToCreds(cfg)
	return creds, nil
}

func (h *IniHandler) mapCfgToCreds(cfg *ini.File) map[string]core.Cred {
	creds := make(map[string]core.Cred)

	for _, sec := range cfg.Sections() {
		// TODO: investigate "DEFAULT" section.
		if sec.Name() == ini.DefaultSection {
			continue
		}

		profile := sec.Name()
		cred := core.Cred{}

		for _, key := range sec.Keys() {
			switch key.Name() {
			case keyAwsAccessKeyID:
				cred.AccessKeyID = sec.Key(keyAwsAccessKeyID).String()
			case keyAwsSecretAccessID:
				cred.SecretAccessKey = sec.Key(keyAwsSecretAccessID).String()
			case keyAwsSessionToken:
				cred.SessionToken = sec.Key(keyAwsSessionToken).String()
			}
		}
		creds[profile] = cred
	}
	return creds
}

// Write overwrite credentials on the file.
func (h *IniHandler) Write(creds map[string]core.Cred) error {
	if h.readOnly {
		return ErrReadOnly
	}

	cfg := h.mapCredsToCfg(creds)
	return cfg.SaveTo(h.filepath)
}

func (h *IniHandler) mapCredsToCfg(creds map[string]core.Cred) *ini.File {
	cfg := ini.Empty()
	
	for profile, cred := range creds {
		sec, _ := cfg.NewSection(profile)

		sec.Key(keyAwsAccessKeyID).SetValue(cred.AccessKeyID)
		sec.Key(keyAwsSecretAccessID).SetValue(cred.SecretAccessKey)
		if cred.SessionToken != "" {
			sec.Key(keyAwsSessionToken).SetValue(cred.SessionToken)
		}
	}
	return cfg
}

// Remove delete the file.
func (h *IniHandler) Remove() error {
	if h.readOnly {
		return ErrReadOnly
	}

	err := os.Remove(h.filepath)
	return err
}
