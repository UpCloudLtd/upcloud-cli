package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecret(t *testing.T) {
	assert.Equal(t, HiddenPlaceholder, Secret("s3cret", false))
	assert.Equal(t, "s3cret", Secret("s3cret", true))
	assert.Equal(t, "", Secret("", false))
}

func TestSecretURI(t *testing.T) {
	uri := "postgres://upadmin:p%40ss@my-pg.example.com:11550/defaultdb?sslmode=require" //nolint:gosec // G101: test fixture, not a real credential
	assert.Equal(t, "postgres://upadmin:xxxxx@my-pg.example.com:11550/defaultdb?sslmode=require", SecretURI(uri, false))
	assert.Equal(t, uri, SecretURI(uri, true))
	assert.Equal(t, "rediss://my-cache.example.com:11550", SecretURI("rediss://my-cache.example.com:11550", false))
	assert.Equal(t, HiddenPlaceholder, SecretURI("postgres://user:pa ss@%zz", false))
	assert.Equal(t, "", SecretURI("", false))
}
