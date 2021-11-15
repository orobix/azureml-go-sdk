package aml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClientEmptyConfig(t *testing.T) {
	a := assert.New(t)

	client, err := NewClient(ClientConfig{}, false)

	a.NotNil(err)
	a.Empty(client)
}

func TestNewClientInvalidAuth(t *testing.T) {
	a := assert.New(t)

	config := ClientConfig{
		ClientId:     "invalid",
		ClientSecret: "invalid",
		TenantId:     "invalid",
	}
	client, err := NewClient(config, false)

	a.Nil(err)
	a.NotEmpty(client)
}
