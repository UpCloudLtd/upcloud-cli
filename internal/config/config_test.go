package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
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

func TestConfig_GetVersion(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		version  string
	}{
		{
			name:     "Removes v prefix",
			version:  "v1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "Removes v prefix with suffix",
			version:  "v1.2.3-rc1",
			expected: "1.2.3-rc1",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			Version = test.version
			actual := GetVersion()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestConfig_LoadKeyring(t *testing.T) {
	// Note that configs defined in environment variables will override configs defined in the config file. Thus, this test will fail if credentials are currently defined as environment variables.
	cfg := New()
	tmpFile, err := os.CreateTemp(os.TempDir(), "")
	assert.NoError(t, err)
	_, err = tmpFile.WriteString("username: unittest")
	assert.NoError(t, err)

	err = keyring.Set("UpCloud", "unittest", "unittest_password")
	assert.NoError(t, err)

	cfg.GlobalFlags.ConfigFile = tmpFile.Name()
	err = cfg.Load()
	assert.NoError(t, err)
	assert.Equal(t, cfg.GetString("username"), "unittest")
	assert.Equal(t, "unittest_password", cfg.GetString("password"))
	t.Cleanup(func() {
		// remove test user from keyring
		assert.NoError(t, keyring.Delete("UpCloud", "unittest"))
	})
}
