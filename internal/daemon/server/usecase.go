package server

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/hanjunlee/awscred/core"
)

type (
	// FileWatcher send a event when the file has a operation.
	FileWatcher interface {
		Watch(context.Context, chan<- fsnotify.Event)
	}

	// CredFileHandler is the manager read and write a credential file.
	CredFileHandler interface {
		Read() (map[string]core.Cred, error)
		Write(map[string]core.Cred) error
	}

	// ConfigFileHandler is the manager read and write a config file.
	ConfigFileHandler interface {
		Read() (map[string]core.Config, error)
		Write(map[string]core.Config) error
	}
)

func mapConfigToCred(c core.Config) core.Cred {
	return core.Cred{
		AccessKeyID:     c.Cache.AccessKeyID,
		SecretAccessKey: c.Cache.SecretAccessKey,
		SessionToken:    c.Cache.SessionToken,
	}
}
