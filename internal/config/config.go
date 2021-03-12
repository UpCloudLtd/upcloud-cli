package config

import (
	// "fmt"
	// "os"
	// "path/filepath"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/UpCloudLtd/cli/internal/terminal"
)

const (
	EnvPrefix             = "UPCLOUD"
	ConfigFileDefaultName = "upctl"
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
	// AppConfig shared config
	AppConfig *Config
)

// New generates a new instance of application config
func New() *Config {
	return &Config{
		Viper: viper.New(),
	}
}

// Config holds the configuration for running upctl
type Config struct {
	Service *service.Service
	Viper   *viper.Viper
	Flags   *pflag.FlagSet
}

// Load() load config to Viper instance
// func (s *Config) Load() error {
// 	s.Viper().SetEnvPrefix(envPrefix)
// 	s.Viper().SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
// 	s.Viper().AutomaticEnv()
// 	s.Viper().SetConfigName("upctl")
// 	s.Viper().SetConfigType("yaml")

// 	if configFile := s.GetString("config"); configFile != "" {
// 		fmt.Printf("1config file set %v\n", configFile)
// 		s.Viper().SetConfigFile(configFile)
// 		s.Viper().SetConfigName(path.Base(configFile))
// 	}

// 	s.Viper().AddConfigPath("$HOME/.config")

// 	if err := s.Viper().ReadInConfig(); err != nil {
// 		return err
// 	}
// 	s.Viper().Set("config", s.Viper().ConfigFileUsed())

// 	return nil

// }

// Viper returns a reference to the viper instance
// func (s *Config) Viper() *viper.Viper {
// 	return s.viper
// }

// // SetNamespace sets the configuration namespace
// func (s *Config) SetNamespace(ns string) {
// 	s.ns = ns
// }

// IsSet return true if the key is set in the current namespace
func (s *Config) IsSet(key string) bool {
	return s.Viper.IsSet(key)
}

// Get return the value of the key in the current namespace
func (s *Config) Get(key string) interface{} {
	return s.Viper.Get(key)
}

// GetString is a convenience method of getting a configuration value in the current namespace as a string
func (s *Config) GetString(key string) string {
	return s.Viper.GetString(key)
}

// // FlagByKey returns pflag.Flag associated with a key in config
// func (s *Config) FlagByKey(key string) *pflag.Flag {
// 	if s.flagSet == nil {
// 		s.flagSet = &pflag.FlagSet{}
// 	}
// 	return s.flagSet.Lookup(key)
// }

// // BoundFlags returns the list of all the flags given to the config
// func (s *Config) BoundFlags() []*pflag.Flag {
// 	if s.flagSet == nil {
// 		s.flagSet = &pflag.FlagSet{}
// 	}
// 	var r []*pflag.Flag
// 	s.flagSet.VisitAll(func(flag *pflag.Flag) {
// 		r = append(r, flag)
// 	})
// 	return r
// }

// // ConfigBindFlagSet sets the config flag set and binds them to the viper instance
// func (s *Config) ConfigBindFlagSet(flags *pflag.FlagSet) {
// 	if flags == nil {
// 		panic("nil flag")
// 	}
// 	if s.flagSet == nil {
// 		s.flagSet = &pflag.FlagSet{}
// 	}
// 	flags.VisitAll(func(flag *pflag.Flag) {
// 		if s.flagSet.Lookup(flag.Name) != nil {
// 			panic("flag exists")
// 		}
// 		_ = s.viper.BindPFlag(s.prependNs(flag.Name), flag)
// 		s.flagSet.AddFlag(flag)
// 	})
// }

// Commonly used keys as accessors

// Output is a convenience method for getting the user specified output
func (s *Config) Output() string {
	return s.Viper.GetString("output")
}

// OutputHuman is a convenience method that returns true if the user specified human-readable output
func (s *Config) OutputHuman() bool {
	return s.Output() == ValueOutputHuman
}

// ClientTimeout is a convenience method that returns the user specified client timeout
func (s *Config) ClientTimeout() time.Duration {
	return s.Viper.GetDuration(KeyClientTimeout)
}

// InteractiveUI is a convenience method that returns true if the user has requested human output and the terminal supports it.
func (s *Config) InteractiveUI() bool {
	return terminal.IsStdoutTerminal() && s.OutputHuman()
}
