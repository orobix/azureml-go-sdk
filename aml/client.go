package aml

import (
	"context"
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"go.uber.org/zap"
)

type Client struct {
	MsalClient confidential.Client
	logger     *zap.SugaredLogger
}

func (c Client) getJwt() (string, error) {
	scopes := []string{DefaultAmlOauthScope}
	c.logger.Debug("Using cached JWT silently...")
	authResult, err := c.MsalClient.AcquireTokenSilent(context.Background(), scopes)
	if err != nil {
		c.logger.Debug("Could not acquire JWT silently, now acquiring it with Client Credential flow...")
		authResult, err = c.MsalClient.AcquireTokenByCredential(context.Background(), scopes)
		c.logger.Debug("JWT acquired")
	}
	return authResult.AccessToken, err
}

func NewClient(clientId string, clientSecret string, tenantId string) (*Client, error) {
	credential, err := confidential.NewCredFromSecret(clientSecret)
	if err != nil {
		return &Client{}, err
	}

	authority := fmt.Sprintf("https://login.microsoftonline.com/%s", tenantId)
	client, err := confidential.New(clientId, credential, confidential.WithAuthority(authority))
	if err != nil {
		return &Client{}, err
	}

	logger, err := zap.NewDevelopment()
	return &Client{MsalClient: client, logger: logger.Sugar()}, nil
}
