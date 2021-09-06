package config_test

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/UpCloudLtd/upcloud-cli/internal/config"
)

func TestSetBoolFlag_EnableDisableFlags(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name          string
		args          []string
		expectedState config.OptionalBoolean
		expectedError string
	}{
		{
			name:          "set to true",
			args:          []string{"--enable-test"},
			expectedState: config.True,
		},
		{
			name:          "set to false",
			args:          []string{"--disable-test"},
			expectedState: config.False,
		},
		{
			name:          "no flag",
			args:          []string{},
			expectedState: config.Unset,
		},
		{
			name:          "both options",
			args:          []string{"--enable-test", "--disable-test"},
			expectedError: "invalid argument \"false\" for \"--disable-test\" flag: cannot set twice",
		},
	} {
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fs := &pflag.FlagSet{}
			var target config.OptionalBoolean
			config.AddEnableDisableFlags(fs, &target, "test", "testing")
			err := fs.Parse(test.args)
			if test.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedState, target)
			} else {
				assert.EqualError(t, err, test.expectedError)
			}
		})
	}
}

func TestSetBoolFlag_EnableOrDisableFlag(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name          string
		defaultValue  bool
		args          []string
		expectedState config.OptionalBoolean
		expectedError string
	}{
		{
			name:          "set to true",
			defaultValue:  false,
			args:          []string{"--enable-test"},
			expectedState: config.True,
		},
		{
			name:          "set to false",
			defaultValue:  true,
			args:          []string{"--disable-test"},
			expectedState: config.False,
		},
		{
			name:          "true stays true with no flag",
			defaultValue:  true,
			args:          []string{},
			expectedState: config.DefaultTrue,
		},
		{
			name:          "false stays false with no flag",
			defaultValue:  false,
			args:          []string{},
			expectedState: config.DefaultFalse,
		},
		{
			name:          "both options passed in when true",
			defaultValue:  true,
			args:          []string{"--enable-test", "--disable-test"},
			expectedError: "unknown flag: --enable-test",
		},
		{
			name:          "both options passed in when false",
			defaultValue:  false,
			args:          []string{"--enable-test", "--disable-test"},
			expectedError: "unknown flag: --disable-test",
		},
	} {
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fs := &pflag.FlagSet{}
			target := config.Unset
			config.AddEnableOrDisableFlag(fs, &target, test.defaultValue, "test", "testing")
			err := fs.Parse(test.args)
			if test.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedState, target)
			} else {
				assert.EqualError(t, err, test.expectedError)
			}
		})
	}
}
