package cloudevent_poller

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Wrap the AssumeRole credentials in a credentials Provider interface
//
// TODO: this feels awkward .. isn't there an easier way?

func NewAssumeRoleCredentialsProvider(credentials *sts.Credentials) *AssumeRoleCredentialsProvider {
	return &AssumeRoleCredentialsProvider{
		AssumeRoleCredentials: credentials,
	}
}

type AssumeRoleCredentialsProvider struct {
	AssumeRoleCredentials *sts.Credentials
}

func (c AssumeRoleCredentialsProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     *c.AssumeRoleCredentials.AccessKeyId,
		SecretAccessKey: *c.AssumeRoleCredentials.SecretAccessKey,
		SessionToken:    *c.AssumeRoleCredentials.SessionToken,
		ProviderName:    "AssumeRoleCredentialsProvider",
	}, nil

}

func (c AssumeRoleCredentialsProvider) IsExpired() bool {
	return c.AssumeRoleCredentials.Expiration.After(time.Now())

}
