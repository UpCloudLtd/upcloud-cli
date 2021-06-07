package config

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestSetBoolFlag_EnableDisableFlags(t *testing.T) {
	for _, test := range []struct {
		name          string
		args          []string
		expectedState OptionalBoolean
		expectedError string
	}{
		{
			name:          "set to true",
			args:          []string{"--enable-test"},
			expectedState: True,
		},
		{
			name:          "set to false",
			args:          []string{"--disable-test"},
			expectedState: False,
		},
		{
			name:          "no flag",
			args:          []string{},
			expectedState: Unset,
		},
		{
			name:          "both options",
			args:          []string{"--enable-test", "--disable-test"},
			expectedError: "invalid argument \"false\" for \"--disable-test\" flag: cannot set twice",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			fs := &pflag.FlagSet{}
			var target OptionalBoolean
			AddEnableDisableFlags(fs, &target, "test", "testing")
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
	for _, test := range []struct {
		name          string
		defaultValue  bool
		args          []string
		expectedState OptionalBoolean
		expectedError string
	}{
		{
			name:          "set to true",
			defaultValue:  false,
			args:          []string{"--enable-test"},
			expectedState: True,
		},
		{
			name:          "set to false",
			defaultValue:  true,
			args:          []string{"--disable-test"},
			expectedState: False,
		},
		{
			name:          "true stays true with no flag",
			defaultValue:  true,
			args:          []string{},
			expectedState: DefaultTrue,
		},
		{
			name:          "false stays false with no flag",
			defaultValue:  false,
			args:          []string{},
			expectedState: DefaultFalse,
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
		t.Run(test.name, func(t *testing.T) {
			fs := &pflag.FlagSet{}
			target := Unset
			AddEnableOrDisableFlag(fs, &target, test.defaultValue, "test", "testing")
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
