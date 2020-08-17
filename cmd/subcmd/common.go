package subcmd

import (
	log "github.com/sirupsen/logrus"
)

func setDebugMode() {
	log.SetLevel(log.DebugLevel)
}
