package domain

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type (
	// Service -
	Service interface {
		GetSecureCredential(*Cred, *Config, string) (*Config, error)
	}

	// CredService -
	CredService struct{}
)

// GetSecureCredential return a new profile with the secure credential.
func (s *CredService) GetSecureCredential(cred *Cred, c *Config, token string) (*Config, error) {
	sess := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cred.AccessKeyID, cred.SecretAccessKey, ""),
	})

	svc := sts.New(sess)
	res, err := svc.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(c.DurationSecond),
		SerialNumber:    aws.String(c.SerialNumber),
		TokenCode:       aws.String(token),
	})
	if err != nil {
		return nil, err
	}

	return &Config{
		DurationSecond: c.DurationSecond,
		SerialNumber:   c.SerialNumber,

		Cache: Cache{
			AccessKeyID:     aws.StringValue(res.Credentials.AccessKeyId),
			SecretAccessKey: aws.StringValue(res.Credentials.SecretAccessKey),
			SessionToken:    aws.StringValue(res.Credentials.SessionToken),
			Expiration:      aws.TimeValue(res.Credentials.Expiration),
		},
	}, err
}
