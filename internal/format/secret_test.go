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
	assert.Equal(t, HiddenPlaceholder, SecretURI("postgres://upadmin@host/db?password=s3cret", false, "s3cret"))
	assert.Equal(t, HiddenPlaceholder, SecretURI("postgres:upadmin:s3cret@host:5432", false, "s3cret"))
	assert.Equal(t, HiddenPlaceholder, SecretURI("postgres://upadmin@host/db?password=p%40ss", false, "p@ss"))
	assert.Equal(t, "postgres://upadmin:xxxxx@host/db", SecretURI("postgres://upadmin:s3cret@host/db", false, "s3cret"))
}
