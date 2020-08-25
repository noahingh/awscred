package server

import (
	"context"
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hanjunlee/awscred/core"
	"github.com/hanjunlee/awscred/internal/daemon/pkg/configfile"
	"github.com/hanjunlee/awscred/internal/daemon/pkg/credfile"
	"github.com/hanjunlee/awscred/internal/daemon/pkg/credwatcher"
	"github.com/hanjunlee/awscred/internal/daemon/pkg/sts"
	"github.com/sirupsen/logrus"
)

type (
	// Interactor manage the credential file and the config file.
	Interactor struct {
		ch              chan fsnotify.Event
		stGenerator     SessionTokenGenerator
		watcher         FileWatcher
		origCredHandler CredFileHandler
		credHandler     CredFileHandler
		confHandler     ConfigFileHandler
		log             *logrus.Entry
	}

	// ProfileInfo is the information of profile.
	ProfileInfo struct {
		Name     string
		On       bool
		Serial   string
		Duration int64
		Expired  string
	}
)

// NewInteractor create a new interactor.
func NewInteractor(origCredPath, credPath, confPath string) *Interactor {
	ch := make(chan fsnotify.Event, 10)

	return &Interactor{
		ch:              ch,
		stGenerator:     sts.NewStsGenerator(),
		watcher:         credwatcher.NewService(origCredPath),
		origCredHandler: credfile.NewIniHandler(true, origCredPath),
		credHandler:     credfile.NewIniHandler(false, credPath),
		confHandler:     configfile.NewIniHandler(confPath),
		log:             logrus.NewEntry(logrus.New()),
	}
}

// StartWatch start to watch the orignal credential file and
// reflect changes into the credential file.
func (i *Interactor) StartWatch(ctx context.Context) {
	i.log.Info("start to watch the aws credentials.")
	go i.watcher.Watch(ctx, i.ch)

	i.log.Debug("run a worker.")
	go i.runWorker(ctx)
}

func (i *Interactor) runWorker(ctx context.Context) {
	for {
		event, more := <-i.ch
		if !more {
			break
		}

		if event.Op == fsnotify.Remove || event.Op == fsnotify.Rename {
			i.log.Warn("the file is disappeared. if the editor use atomic saves you should restart the daemon(https://github.com/fsnotify/fsnotify/issues/17).")
			continue
		}

		if err := i.reflect(); err != nil {
			i.log.Errorf("failed to reflect: %s", err)
		}
	}
}

// Reflect reflect the original credential file, only for enabled profiles.
func (i *Interactor) Reflect() error {
	i.log.Info("reflect on the awscred credentials.")
	return i.reflect()
}

// reflect the original credential file, only for enabled profiles.
func (i *Interactor) reflect() error {
	origCreds, err := i.origCredHandler.Read()
	if err != nil {
		return fmt.Errorf("failed to read the original credential file: %s", err)
	}

	confs, err := i.confHandler.Read()
	if err != nil {
		return fmt.Errorf("failed to read the awscred config file: %s", err)
	}

	reflected := make(map[string]core.Cred)
	for profile, orig := range origCreds {
		conf, ok := confs[profile]

		if !ok {
			i.log.Debugf("the config of profile doesn't exist, reflect the original credential: \"%s\".", profile)
			reflected[profile] = orig
			continue
		}

		if !conf.On {
			i.log.Debugf("the config of profile is disabled: \"%s\".", profile)
			reflected[profile] = orig
			continue
		}

		i.log.Debugf("the config of profile is disabled: \"%s\".", profile)
		reflected[profile] = mapConfigToCred(conf)
	}

	if err := i.credHandler.Write(reflected); err != nil {
		return fmt.Errorf("failed to write the awscred credential file")
	}

	return nil
}

// Terminate stop to watch the original credential file and remove the reflected credential file.
func (i *Interactor) Terminate() error {
	close(i.ch)
	return i.credHandler.Remove()
}

// On set the profile enabled, and if the configuration doesn't exist
// it create a new configuration.
func (i *Interactor) On(profile string) error {
	_, ok, err := i.GetOriginalCred(profile)
	if err != nil {
		return err
	}
	if ok != true {
		return fmt.Errorf("there's no such a profile: %s", profile)
	}

	config, ok, err := i.GetConfig(profile)
	if err != nil {
		return err
	}
	config.On = true

	if err := i.SetConfig(profile, config); err != nil {
		return err
	}

	return nil
}

// Off set the profile disabled.
func (i *Interactor) Off(profile string) error {
	_, ok, err := i.GetOriginalCred(profile)
	if err != nil {
		return err
	}
	if ok != true {
		return fmt.Errorf("there's no such a profile: %s", profile)
	}

	config, ok, err := i.GetConfig(profile)
	if err != nil {
		return err
	}
	config.On = false

	if err := i.SetConfig(profile, config); err != nil {
		return err
	}
	return nil
}

// Gen generate the secure token from STS.
func (i *Interactor) Gen(profile, token string) error {
	cred, ok, err := i.GetOriginalCred(profile)
	if err != nil {
		return err
	}
	if ok != true {
		return fmt.Errorf("there's no such a profile: %s", profile)
	}

	config, ok, err := i.GetConfig(profile)
	if err != nil {
		return err
	}
	if !config.On {
		return fmt.Errorf("it's disabled, set the profile enabled: %s", profile)
	}

	sc, err := i.stGenerator.Generate(cred, config, token)
	if err != nil {
		return fmt.Errorf("failed to get the secure credential from STS: %s", err)
	}

	// cache secure credentail.
	config.Cache = sc
	if err := i.SetConfig(profile, config); err != nil {
		return fmt.Errorf("failed to set the profile")
	}

	if err := i.SetCred(profile, mapConfigToCred(config)); err != nil {
		return fmt.Errorf("failed to set the profile")
	}

	return nil
}

// GetOriginalCred get the credential of profile.
func (i *Interactor) GetOriginalCred(profile string) (core.Cred, bool, error) {
	creds, err := i.origCredHandler.Read()
	if err != nil {
		return core.Cred{}, false, err
	}

	cred, ok := creds[profile]
	if !ok {
		return core.Cred{}, false, nil
	}

	return cred, ok, nil
}

// SetCred set a new credential for the profile.
func (i *Interactor) SetCred(profile string, cred core.Cred) error {
	creds, err := i.credHandler.Read()
	if err != nil {
		return err
	}

	creds[profile] = cred

	if err := i.credHandler.Write(creds); err != nil {
		return err
	}

	return nil
}

// GetConfig get the config of profile.
func (i *Interactor) GetConfig(profile string) (core.Config, bool, error) {
	configs, err := i.confHandler.Read()
	if err != nil {
		return core.Config{}, false, err
	}

	config, ok := configs[profile]
	if !ok {
		return core.Config{}, false, nil
	}

	return config, ok, nil
}

// SetConfig set the config of profile.
func (i *Interactor) SetConfig(profile string, conf core.Config) error {
	configs, err := i.confHandler.Read()
	if err != nil {
		return err
	}

	configs[profile] = conf

	if err := i.confHandler.Write(configs); err != nil {
		return err
	}

	return nil
}

// GetProfileList return profiles.
func (i *Interactor) GetProfileList() ([]ProfileInfo, error) {
	const (
		NA = "N/A"
	)
	ps := make([]ProfileInfo, 0)

	creds, err := i.origCredHandler.Read()
	if err != nil {
		i.log.Errorf("failed to read the aws credentials: %s", err)
		return ps, err
	}

	configs, err := i.confHandler.Read()
	if err != nil {
		i.log.Errorf("failed to read the awscred configs: %s", err)
		return ps, err
	}

	for profile := range creds {
		var (
			pi ProfileInfo
		)
		conf, ok := configs[profile]
		if !ok {
			pi = ProfileInfo{
				Name: profile,
				On:   false,
			}
			ps = append(ps, pi)
			continue
		}

		pi = ProfileInfo{
			Name:     profile,
			On:       conf.On,
			Serial:   conf.SerialNumber,
			Duration: conf.DurationSecond,
			Expired:  conf.Cache.Expiration.Format(time.RFC3339),
		}
		ps = append(ps, pi)
	}

	return ps, nil
}
