package awsctr

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/juju/errors"
)

const (
	DefaultRegion  = "ap-northeast-1"
	EnvAccessKeyID = "AWS_ACCESS_KEY_ID"
	EnvSecretKey   = "AWS_SECRET_ACCESS_KEY"
)

func NewIAMClient() (*iam.IAM, error) {
	sess, err := NewSessionFromEnv(DefaultRegion)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return iam.New(sess), nil
}

func NewSessionFromEnv(region string) (*session.Session, error) {
	accessKeyID := os.Getenv(EnvAccessKeyID)
	if accessKeyID == "" {
		return nil, errors.Errorf("environment variable %q is empty", EnvAccessKeyID)
	}

	secretKey := os.Getenv(EnvSecretKey)
	if secretKey == "" {
		return nil, errors.Errorf("environment variable %q is empty", EnvSecretKey)
	}

	return session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretKey, ""),
	})
}
