package core

import "time"

type (
	// SessionToken -
	SessionToken struct {
		AccessKeyID     string
		SecretAccessKey string
		SessionToken    string
		Expiration      time.Time
	}

	// Config -
	Config struct {
		On             bool
		SerialNumber   string
		DurationSecond int64

		Cache SessionToken
	}
)
