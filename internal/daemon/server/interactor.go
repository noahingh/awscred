package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/hanjunlee/awsmonkey/core"
	"github.com/sirupsen/logrus"
)

type (
	// Interactor manage the credential file and the config file.
	Interactor struct {
		ch              chan fsnotify.Event
		watcher         FileWatcher
		origCredHandler CredFileHandler
		credHandler     CredFileHandler
		confHandler     ConfigFileHandler
		log             *logrus.Entry
	}
)

// StartWatch start to watch the orignal credential file and
// reflect changes into the credential file.
func (i *Interactor) StartWatch(ctx context.Context) {
	i.log.Info("start to watch.")
	go i.watcher.Watch(ctx, i.ch)

	i.log.Info("run a worker.")
	go i.runWorker(ctx)
}

func (i *Interactor) runWorker(ctx context.Context) {
	for {
		event, more := <-i.ch
		if !more {
			break
		}

		if event.Op == fsnotify.Remove || event.Op == fsnotify.Rename {
			i.log.Warn("the file is disappeared.")
			continue
		}

		if err := i.reflect(); err != nil {
			i.log.Error("failed to reflect: %s", err)
		}
	}
}

// reflect the original credential file, only for enabled profiles.
func (i *Interactor) reflect() error {
	origCreds, err := i.origCredHandler.Read()
	if err != nil {
		return fmt.Errorf("failed to read the original credential file: %s", err)
	}

	confs, err := i.confHandler.Read()
	if err != nil {
		return fmt.Errorf("failed to read the awsmockey config file: %s", err)
	}

	reflected := make(map[string]core.Cred)
	for profile, orig := range origCreds {
		conf, ok := confs[profile]

		if !ok {
			reflected[profile] = orig
			continue
		}

		if t, err := strconv.ParseBool(conf.On); err != nil || !t {
			reflected[profile] = orig
		}

		reflected[profile] = mapConfigToCred(conf)
	}

	if err := i.credHandler.Write(reflected); err != nil {
		return fmt.Errorf("failed to write the awsmockey credential file")
	}

	return nil
}
