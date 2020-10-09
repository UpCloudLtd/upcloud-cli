package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/all"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/terminal"
	"github.com/UpCloudLtd/cli/internal/validation"
)

const envPrefix = "UPCLOUD"

var (
	mainConfig = config.New(viper.New())
	mc         = commands.BuildCommand(
		&mainCommand{BaseCommand: commands.New("upctl", "UpCloud command line client")},
		nil, mainConfig)
)

type completionCommand struct {
	*commands.BaseCommand
}

func (s *completionCommand) MakeExecuteCommand() func(args []string) error {
	return func(args []string) error {
		switch args[0] {
		case "bash":
			_ = s.Cobra().Root().GenBashCompletion(os.Stdout)
		}
		return nil
	}
}

type mainCommand struct {
	*commands.BaseCommand
}

func (s *mainCommand) initConfig(outputErrors bool) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil && outputErrors {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	s.Config().Viper().SetEnvPrefix(envPrefix)
	s.Config().Viper().SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	s.Config().Viper().AutomaticEnv()
	s.Config().Viper().SetConfigName(".upctl")
	s.Config().Viper().SetConfigType("yaml")

	if configFile := s.Config().GetString("config"); configFile != "" {
		s.Config().Viper().SetConfigFile(configFile)
		s.Config().Viper().SetConfigName(path.Base(configFile))
		configDir := path.Dir(configFile)
		if configDir != "." && configDir != dir {
			viper.AddConfigPath(configDir)
		}
	}

	s.Config().Viper().AddConfigPath(dir)
	s.Config().Viper().AddConfigPath(".")
	s.Config().Viper().AddConfigPath("$HOME")

	if err := s.Config().Viper().ReadInConfig(); err != nil && outputErrors {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: config file load error: %v\n", err)
	}
	s.Config().Viper().Set("config", s.Config().Viper().ConfigFileUsed())
}

func (s *mainCommand) InitCommand() {
	s.Cobra().SilenceErrors = true
	s.Cobra().SilenceUsage = true
	flags := &pflag.FlagSet{}
	flags.String("config", "", "Config file")
	flags.String("output", "human", "Output format (supported: json, yaml and human")
	flags.String("colours", "auto", "Use terminal colours (supported: auto, true, false)")
	flags.Duration("client-timeout", 600*time.Second, "Timeout for requests")
	s.AddPersistentFlags(flags)

	commands.BuildCommand(
		&completionCommand{commands.New("completion", "Generate shell completion code")},
		s, config.New(mainConfig.Viper()))

	s.SetConfigLoader(func(config *config.Config, loadContext int) {
		s.initConfig(loadContext == commands.ConfigLoadContextHelp)
	})
	all.BuildCommands(s, s.Config())
}

func (s *mainCommand) MakePersistentPreExecuteCommand() func(args []string) error {
	return func(args []string) error {
		if err := validation.Value(s.Config().GetString("output"), "json", "yaml", "human"); err != nil {
			return fmt.Errorf("invalid output: %v", err)
		}
		if err := validation.Value(s.Config().GetString("colours"), "auto", "true", "false", "on", "off"); err != nil {
			return fmt.Errorf("invalid colours: %v", err)
		}
		switch strings.ToLower(s.Config().GetString("colours")) {
		case "on", "true":
			terminal.ForceColours(true)
		case "off", "false":
			terminal.ForceColours(false)
		}
		return nil
	}
}

func main() {
	if err := mc.Cobra().Execute(); err != nil {
		mc.HandleError(err)
		os.Exit(1)
	}
}
