package config

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/UpCloudLtd/cli/internal/terminal"
)

const (
	// KeyClientTimeout defines the viper configuration key used to define client timeout
	KeyClientTimeout = "client-timeout"
	// KeyOutput defines the viper configuration key used to define the output
	KeyOutput = "output"
	// ValueOutputHuman defines the viper configuration value used to define human-readable output
	ValueOutputHuman = "human"
	// ValueOutputYAML defines the viper configuration value used to define YAML output
	ValueOutputYAML = "yaml"
	// ValueOutputJSON defines the viper configuration value used to define JSON output
	ValueOutputJSON = "json"
)

var (
	// Version contains the current version.
	Version = "dev"
	// BuildDate contains a string with the build date.
	BuildDate = "unknown"
)

// New returns a new instance of Config bound to the given viper instance
func New(viper *viper.Viper) *Config {
	return &Config{viper: viper}
}

// Config holds the configuration for running upctl
type Config struct {
	viper   *viper.Viper
	ns      string
	flagSet *pflag.FlagSet
}

// Viper returns a reference to the viper instance
func (s *Config) Viper() *viper.Viper {
	return s.viper
}

// SetNamespace sets the configuration namespace
func (s *Config) SetNamespace(ns string) {
	s.ns = ns
}

// IsSet return true if the key is set in the current namespace
func (s *Config) IsSet(key string) bool {
	return s.viper.IsSet(s.prependNs(key))
}

// Get return the value of the key in the current namespace
func (s *Config) Get(key string) interface{} {
	return s.viper.Get(s.prependNs(key))
}

// GetString is a convenience method of getting a configuration value in the current namespace as a string
func (s *Config) GetString(key string) string {
	return s.viper.GetString(s.prependNs(key))
}

// Top returns a *copy* of the Config with no namespace
func (s *Config) Top() *Config {
	clone := *s
	clone.SetNamespace("")
	return &clone
}

// FlagByKey returns pflag.Flag associated with a key in config
func (s *Config) FlagByKey(key string) *pflag.Flag {
	if s.flagSet == nil {
		s.flagSet = &pflag.FlagSet{}
	}
	return s.flagSet.Lookup(key)
}

// BoundFlags returns the list of all the flags given to the config
func (s *Config) BoundFlags() []*pflag.Flag {
	if s.flagSet == nil {
		s.flagSet = &pflag.FlagSet{}
	}
	var r []*pflag.Flag
	s.flagSet.VisitAll(func(flag *pflag.Flag) {
		r = append(r, flag)
	})
	return r
}

// ConfigBindFlagSet sets the config flag set and binds them to the viper instance
func (s *Config) ConfigBindFlagSet(flags *pflag.FlagSet) {
	if flags == nil {
		return
	}
	if s.flagSet == nil {
		s.flagSet = &pflag.FlagSet{}
	}
	flags.VisitAll(func(flag *pflag.Flag) {
		if s.flagSet.Lookup(flag.Name) != nil {
			return
		}
		_ = s.viper.BindPFlag(s.prependNs(flag.Name), flag)
		s.flagSet.AddFlag(flag)
	})
}

func (s *Config) prependNs(key string) string {
	if s.ns == "" {
		return key
	}
	return fmt.Sprintf("%s.%s", s.ns, key)
}

// Commonly used keys as accessors

// Output is a convenience method for getting the user specified output
func (s *Config) Output() string {
	return s.viper.GetString(KeyOutput)
}

// OutputHuman is a convenience method that returns true if the user specified human-readable output
func (s *Config) OutputHuman() bool {
	return s.Output() == ValueOutputHuman
}

// ClientTimeout is a convenience method that returns the user specified client timeout
func (s *Config) ClientTimeout() time.Duration {
	return s.viper.GetDuration(KeyClientTimeout)
}

// InteractiveUI is a convenience method that returns true if the user has requested human output and the terminal supports it.
func (s *Config) InteractiveUI() bool {
	return terminal.IsStdoutTerminal() && s.OutputHuman()
}
