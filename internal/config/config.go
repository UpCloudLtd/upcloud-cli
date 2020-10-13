package config

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/UpCloudLtd/cli/internal/terminal"
)

const (
	ConfigKeyOutput        = "output"
	ConfigValueOutputHuman = "human"
	ConfigValueOuputYaml   = "yaml"
	ConfigValueOuputJson   = "json"
	ConfigKeyClientTimeout = "client-timeout"
)

func New(viper *viper.Viper) *Config {
	return &Config{viper: viper}
}

type Config struct {
	viper   *viper.Viper
	ns      string
	flagSet *pflag.FlagSet
}

func (s *Config) Viper() *viper.Viper {
	return s.viper
}

func (s *Config) SetNamespace(ns string) {
	s.ns = ns
}

func (s *Config) IsSet(key string) bool {
	return s.viper.IsSet(s.prependNs(key))
}

func (s *Config) Get(key string) interface{} {
	return s.viper.Get(s.prependNs(key))
}

func (s *Config) GetString(key string) string {
	return s.viper.GetString(s.prependNs(key))
}

func (s *Config) Top() *Config {
	clone := *s
	clone.SetNamespace("")
	return &clone
}

func (s *Config) FlagByKey(key string) *pflag.Flag {
	if s.flagSet == nil {
		s.flagSet = &pflag.FlagSet{}
	}
	return s.flagSet.Lookup(key)
}

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

func (s *Config) ConfigBindFlagSet(flags *pflag.FlagSet) {
	if flags == nil {
		panic("Nil flagset")
	}
	if s.flagSet == nil {
		s.flagSet = &pflag.FlagSet{}
	}
	flags.VisitAll(func(flag *pflag.Flag) {
		if s.flagSet.Lookup(flag.Name) != nil {
			panic(fmt.Sprintf("key %s already bound", flag.Name))
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

func (s *Config) Output() string {
	return s.viper.GetString(ConfigKeyOutput)
}

func (s *Config) OutputHuman() bool {
	return s.Output() == ConfigValueOutputHuman
}

func (s *Config) ClientTimeout() time.Duration {
	return s.viper.GetDuration(ConfigKeyClientTimeout)
}

func (s *Config) InteractiveUI() bool {
	return terminal.IsStdoutTerminal() && s.OutputHuman()
}
