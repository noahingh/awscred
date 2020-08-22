package cmd

import (
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

const (
	defaultPort = 5126
)

var (
	homeDir, _ = homedir.Dir()
)

func setDebugMode() {
	log.SetLevel(log.DebugLevel)
}
