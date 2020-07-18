package sts

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hanjunlee/awscred/core"
)

type (
	// Generator -
	Generator struct {
	}
)

// GetSecureCredential return a new profile with the secure credential.
func (g *Generator) GetSecureCredential(cred core.Cred, conf core.Config, token string) (core.SessionToken, error) {
	sess := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cred.AccessKeyID, cred.SecretAccessKey, ""),
	})

	res, err := sts.New(sess).GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(conf.DurationSecond),
		SerialNumber:    aws.String(conf.SerialNumber),
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
