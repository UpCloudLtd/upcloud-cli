package config

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/clierrors"
	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/zalando/go-keyring"

	"github.com/UpCloudLtd/upcloud-go-api/credentials"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
	"github.com/adrg/xdg"
	"github.com/gemalto/flume"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

	// env vars custom prefix
	envPrefix = "UPCLOUD"
)

var (
	// Version contains the current version.
	Version = "dev"
	// BuildDate contains a string with the build date.
	BuildDate = "unknown"
	// flume logger for config, that will be passed to log pckg
	logger = flume.New("config")
)

// New returns a new instance of Config bound to the given viper instance
func New() *Config {
	ctx, cancel := context.WithCancel(context.Background())
	return &Config{viper: viper.New(), context: ctx, cancel: cancel}
}

// GlobalFlags holds information on the flags shared among all commands
type GlobalFlags struct {
	ConfigFile    string        `valid:"-"`
	ClientTimeout time.Duration `valid:"-"`
	Debug         bool          `valid:"-"`
	OutputFormat  string        `valid:"in(human|json|yaml)"`
	NoColours     OptionalBoolean
	ForceColours  OptionalBoolean
}

// Config holds the configuration for running upctl
type Config struct {
	viper       *viper.Viper
	flagSet     *pflag.FlagSet
	cancel      context.CancelFunc
	context     context.Context //nolint: containedctx // This is where the top-level context is stored
	GlobalFlags GlobalFlags
}

// Load loads config and sets up service
func (s *Config) Load() error {
	v := s.Viper()

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()
	v.SetConfigName("upctl")
	v.SetConfigType("yaml")

	configFile := s.GlobalFlags.ConfigFile
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		// Support XDG default config home dir and common config dirs
		v.AddConfigPath(xdg.ConfigHome)
		v.AddConfigPath("$HOME/.config") // for MacOS as XDG config is not common
	}

	// Attempt to read the config file, ignoring only config file not found errors
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("unable to parse config from file '%v': %w", v.ConfigFileUsed(), err)
		}
	}

	creds, err := credentials.Parse(credentials.Credentials{
		Username: v.GetString("username"),
		Password: v.GetString("password"),
		Token:    v.GetString("token"),
	})
	if err == nil {
		v.Set("username", creds.Username)
		v.Set("password", creds.Password)
		v.Set("token", creds.Token)
	}

	v.Set("config", v.ConfigFileUsed())

	settings := v.AllSettings()
	// sanitize password before logging settings
	if _, ok := settings["password"]; ok {
		settings["password"] = "[REDACTED]"
	}
	if _, ok := settings["token"]; ok {
		settings["token"] = "[REDACTED]"
	}
	logger.Debug("viper initialized", "settings", settings)

	return nil
}

// Viper returns a reference to the viper instance
func (s *Config) Viper() *viper.Viper {
	return s.viper
}

// IsSet return true if the key is set in the current namespace
func (s *Config) IsSet(key string) bool {
	return s.viper.IsSet(key)
}

// Get return the value of the key in the current namespace
func (s *Config) Get(key string) interface{} {
	return s.viper.Get(key)
}

// GetString is a convenience method of getting a configuration value in the current namespace as a string
func (s *Config) GetString(key string) string {
	return s.viper.GetString(key)
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
		panic("Nil flagset")
	}
	flags.VisitAll(func(flag *pflag.Flag) {
		_ = s.viper.BindPFlag(flag.Name, flag)
		// s.flagSet.AddFlag(flag)
	})
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

func (s *Config) Cancel() {
	s.cancel()
}

func (s *Config) Context() context.Context {
	return s.context
}

// CreateService creates a new service instance and puts in the conf struct
func (s *Config) CreateService() (internal.AllServices, error) {
	username := s.GetString("username")
	password := s.GetString("password")
	token := s.GetString("token")

	if token == "" && (username == "" || password == "") {
		// This might give silghtly unexpected results on OS X, as xdg.ConfigHome points to ~/Library/Application Support
		// while we really use/prefer/document ~/.config - which does work on osx as well but won't be displayed here.
		configDetails := fmt.Sprintf("default location %s", filepath.Join(xdg.ConfigHome, "upctl.yaml"))
		if s.GetString("config") != "" {
			configDetails = fmt.Sprintf("used %s", s.GetString("config"))
		}
		return nil, clierrors.MissingCredentialsError{ConfigFile: configDetails, ServiceName: credentials.KeyringServiceName}
	}

	configs := []client.ConfigFn{
		client.WithTimeout(s.ClientTimeout()),
	}
	if token != "" {
		configs = append(configs, client.WithBearerAuth(token))
	} else {
		configs = append(configs, client.WithBasicAuth(username, password))
	}

	client := client.New("", "", configs...)
	client.UserAgent = fmt.Sprintf("upctl/%s", GetVersion())

	svc := service.New(client)
	return svc, nil
}

func GetVersion() string {
	version := getVersion()
	re := regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+.*`)
	if re.MatchString(version) {
		return version[1:]
	}
	return version
}

func SaveTokenToKeyring(token string) error {
	return keyring.Set(credentials.KeyringServiceName, credentials.KeyringTokenUser, token)
}

func getVersion() string {
	// Version was overridden during the build
	if Version != "dev" {
		return Version
	}

	// Try to read version from build info
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		version := buildInfo.Main.Version
		if version != "(devel)" && version != "" {
			return version
		}

		settingsMap := make(map[string]string)
		for _, setting := range buildInfo.Settings {
			settingsMap[setting.Key] = setting.Value
		}

		version = "dev"
		if rev, ok := settingsMap["vcs.revision"]; ok {
			version = fmt.Sprintf("%s-%s", version, rev[:8])
		}

		if dirty, ok := settingsMap["vcs.modified"]; ok && dirty == "true" {
			return version + "-dirty"
		}

		return version
	}

	// Fallback to the default value
	return Version
}
