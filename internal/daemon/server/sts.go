package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/hanjunlee/awscred/core"
)

type (
	// SecureTokenGenerator -
	SecureTokenGenerator struct {
		stssvc stsiface.STSAPI
	}
)

// NewCredService return the credential service.
func NewCredService(cred core.Cred) *SecureTokenGenerator {
	sess := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cred.AccessKeyID, cred.SecretAccessKey, ""),
	})

	svc := sts.New(sess)
	return &SecureTokenGenerator{
		stssvc: svc,
	}
}

// GetSecureCredential return a new profile with the secure credential.
func (s *SecureTokenGenerator) GetSecureCredential(c core.Config, token string) (core.SessionToken, error) {
	res, err := s.stssvc.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(c.DurationSecond),
		SerialNumber:    aws.String(c.SerialNumber),
		TokenCode:       aws.String(token),
	})
	if err != nil {
		return core.SessionToken{}, err
	}

	return core.SessionToken{
		AccessKeyID:     aws.StringValue(res.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(res.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(res.Credentials.SessionToken),
		Expiration:      aws.TimeValue(res.Credentials.Expiration),
	}, nil
}
