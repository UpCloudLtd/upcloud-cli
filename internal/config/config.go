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
	ConfigKeyClientTimeout = "client-timeout"
)

func New(viper *viper.Viper) *Config {
	return &Config{viper: viper}
}

type Config struct {
	viper      *viper.Viper
	ns         string
	boundFlags map[string]*pflag.Flag
}

func (s *Config) Viper() *viper.Viper {
	return s.viper
}

func (s *Config) SetNamespace(ns string) {
	s.ns = ns
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

func (s *Config) BoundFlags() map[string]*pflag.Flag {
	return s.boundFlags
}

func (s *Config) ConfigBindFlag(key string, flag *pflag.Flag) {
	if flag == nil {
		panic("Nil flag bound")
	}
	if s.boundFlags == nil {
		s.boundFlags = make(map[string]*pflag.Flag)
	}
	s.boundFlags[key] = flag
	_ = s.viper.BindPFlag(s.prependNs(key), flag)
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
