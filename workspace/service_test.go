package workspace

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClientEmptyConfig(t *testing.T) {
	a := assert.New(t)

	client, err := New(Config{}, false)

	a.NotNil(err)
	a.Empty(client)
}

func TestNewClientInvalidAuth(t *testing.T) {
	a := assert.New(t)

	config := Config{
		ClientId:     "invalid",
		ClientSecret: "invalid",
		TenantId:     "invalid",
	}
	client, err := New(config, false)

	a.Nil(err)
	a.NotEmpty(client)
}
