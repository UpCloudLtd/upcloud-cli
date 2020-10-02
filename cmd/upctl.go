package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/all"
	"github.com/UpCloudLtd/cli/internal/terminal"
	"github.com/UpCloudLtd/cli/internal/validation"
)

const envPrefix = "UPCLOUD"

var (
	mainConfig = viper.New()
	mc         = commands.BuildCommand(&mainCommand{Command: commands.New("upctl", "UpCloud command line client")}, nil, viper.New())
)

type completionCommand struct {
	commands.Command
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
	commands.Command
}

func (s *mainCommand) initConfig(outputErrors bool) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil && outputErrors {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	s.Config().SetEnvPrefix(envPrefix)
	s.Config().AutomaticEnv()
	s.Config().SetConfigName(".upctl")
	s.Config().SetConfigType("yaml")

	if configFile := s.Config().GetString("config"); configFile != "" {
		s.Config().SetConfigFile(configFile)
		s.Config().SetConfigName(path.Base(configFile))
		configDir := path.Dir(configFile)
		if configDir != "." && configDir != dir {
			viper.AddConfigPath(configDir)
		}
	}

	s.Config().AddConfigPath(dir)
	s.Config().AddConfigPath(".")
	s.Config().AddConfigPath("$HOME")

	if err := s.Config().ReadInConfig(); err != nil && outputErrors {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: config file load error: %v\n", err)
	}
	s.Config().Set("config", s.Config().ConfigFileUsed())
}

func (s *mainCommand) InitCommand() {
	s.Cobra().SilenceErrors = true
	s.Cobra().SilenceUsage = true
	s.Cobra().PersistentFlags().String("config", "", "Config file")
	s.ConfigBindFlag("config", s.Cobra().PersistentFlags().Lookup("config"))
	s.Cobra().PersistentFlags().String("output", "human", "Output format (supported: json, yaml and human")
	s.ConfigBindFlag("output", s.Cobra().PersistentFlags().Lookup("output"))
	s.Cobra().PersistentFlags().String("colours", "auto", "Output format (supported: auto, true, false)")
	s.ConfigBindFlag("colours", s.Cobra().PersistentFlags().Lookup("colours"))

	all.BuildCommands(s, s.Config())
	commands.BuildCommand(&completionCommand{commands.New("completion", "Generate shell completion code")}, s, mainConfig)

	s.SetConfigLoader(func(config *viper.Viper, loadContext int) {
		s.initConfig(loadContext == commands.ConfigLoadContextHelp)
	})
}

func (s *mainCommand) MakePersistentPreExecuteCommand() func(args []string) error {
	return func(args []string) error {
		s.initConfig(false)
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
