package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_LoadInvalidYAML(t *testing.T) {
	cfg := New()
	tmpFile, err := os.CreateTemp(os.TempDir(), "")
	assert.NoError(t, err)
	_, err = tmpFile.WriteString("usernamd:sdkfo\npassword: foo")
	assert.NoError(t, err)

	cfg.GlobalFlags.ConfigFile = tmpFile.Name()
	err = cfg.Load()
	assert.EqualError(t, err, fmt.Sprintf("unable to parse config from file '%s': While parsing config: yaml: line 2: mapping values are not allowed in this context", tmpFile.Name()))
}

func TestConfig_Load(t *testing.T) {
	cfg := New()
	tmpFile, err := os.CreateTemp(os.TempDir(), "")
	assert.NoError(t, err)
	_, err = tmpFile.WriteString("username: sdkfo\npassword: foo")
	assert.NoError(t, err)

	cfg.GlobalFlags.ConfigFile = tmpFile.Name()
	err = cfg.Load()
	assert.NoError(t, err)
	assert.NotEmpty(t, cfg.GetString("username"))
	assert.NotEmpty(t, cfg.GetString("password"))
}
