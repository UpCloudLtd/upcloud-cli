package core

import (
	"fmt"
	// "io"
	// "os"

	// "path"
	// "path/filepath"
	// "runtime"
	"strings"

	// "time"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	// "gopkg.in/yaml.v3"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/namespace/root"
	// "github.com/UpCloudLtd/cli/internal/commands/all"
	"github.com/UpCloudLtd/cli/internal/config"
	// "github.com/UpCloudLtd/cli/internal/terminal"
	// "github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/cli/internal/upapi"
	// "github.com/UpCloudLtd/cli/internal/validation"
)

// type completionCommand struct {
// 	*commands.BaseCommand
// }

// func (s *completionCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
// 	return func(args []string) (interface{}, error) {
// 		if len(args) != 1 {
// 			return nil, fmt.Errorf("shell name is requred")
// 		}
// 		shellName := args[0]

// 		if shellName == "bash" {
// 			err := s.Cobra().Root().GenBashCompletion(os.Stdout)
// 			return nil, err
// 		}

// 		return nil, fmt.Errorf("completion for %s is not supported", shellName)
// 	}
// }

// type versionCommand struct {
// 	*commands.BaseCommand
// }

// func (s *versionCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
// 	return func(args []string) (interface{}, error) {
// 		return fmt.Printf(
// 			"Upctl %v\n\tBuild date: %v\n\tBuilt with: %v",
// 			config.Version, config.BuildDate, runtime.Version(),
// 		)
// 	}
// }

// type defaultsCommand struct {
// 	*commands.BaseCommand
// }

// func (s *defaultsCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
// 	return func(args []string) (interface{}, error) {
// 		return nil, s.HandleOutput(os.Stdout, nil)
// 	}
// }

// func (s *defaultsCommand) HandleOutput(io.Writer, interface{}) error {
// 	output := config.ValueOutputYAML
// 	if s.Config().Top().IsSet(config.KeyOutput) && !s.Config().OutputHuman() {
// 		output = s.Config().Output()
// 	}
// 	if output == config.ValueOutputJSON {
// 		return fmt.Errorf("only yaml output is supported for this command")
// 	}

// 	var commandsWithFlags []commands.Command
// nextChild:
// 	for _, cmd := range append([]commands.Command{s.Parent()}, s.Parent().Children()...) {
// 		if len(cmd.Config().BoundFlags()) > 0 {
// 			commandsWithFlags = append(commandsWithFlags, cmd)
// 		} else {
// 			for _, ccmd := range cmd.Children() {
// 				if len(ccmd.Config().BoundFlags()) > 0 {
// 					commandsWithFlags = append(commandsWithFlags, cmd)
// 					continue nextChild
// 				}
// 			}
// 		}
// 	}
// 	topNode := &yaml.Node{Kind: yaml.DocumentNode}
// 	topNode.HeadComment = "UpCtl configuration defaults"
// 	curNode := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
// 	topNode.Content = append(topNode.Content, curNode)
// 	commandToNode := map[commands.Command]*yaml.Node{s.Parent(): curNode}
// 	var prev commands.Command
// 	for _, cmd := range commandsWithFlags {
// 		if prev != nil && prev != cmd.Parent() {
// 			curNode = commandToNode[cmd.Parent()]
// 		}
// 		// Skip top level command name
// 		if cmd.Parent() != nil {
// 			curNode.Content = append(curNode.Content,
// 				&yaml.Node{
// 					Kind:  yaml.ScalarNode,
// 					Tag:   "!!str",
// 					Value: cmd.Use(),
// 				})
// 			curNode.Content = append(curNode.Content, &yaml.Node{
// 				Kind: yaml.MappingNode,
// 				Tag:  "!!map",
// 			})
// 			curNode = curNode.Content[len(curNode.Content)-1]
// 		}
// 		commandToNode[cmd] = curNode
// 		if len(cmd.Config().BoundFlags()) > 0 {
// 			for _, flag := range cmd.Config().BoundFlags() {
// 				keyNode := &yaml.Node{
// 					Kind:        yaml.ScalarNode,
// 					Tag:         "!!str",
// 					HeadComment: flag.Usage,
// 					Value:       flag.Name,
// 					FootComment: "\r",
// 				}
// 				valueNode := &yaml.Node{}
// 				s.applyYamlScalarValue(valueNode, flag.Value.Type(), flag.DefValue)
// 				if valueNode.Kind == 0 &&
// 					(strings.HasSuffix(flag.Value.Type(), "Slice") ||
// 						strings.HasSuffix(flag.Value.Type(), "Array")) {

// 					styp := flag.Value.Type()[0 : len(flag.Value.Type())-5]
// 					valueNode.Kind = yaml.SequenceNode
// 					valueNode.Style = yaml.FlowStyle
// 					valueNode.Tag = "!!seq"
// 					flagSvals := strings.TrimPrefix(flag.Value.String(), "[")
// 					flagSvals = strings.TrimSuffix(flagSvals, "]")
// 					for _, fsv := range strings.Split(flagSvals, ",") {
// 						seqNode := &yaml.Node{}
// 						s.applyYamlScalarValue(seqNode, styp, fsv)
// 						if seqNode.Kind == 0 {
// 							panic(fmt.Sprintf(
// 								"cannot marshal %q value %q type %q",
// 								flag.Name, flag.Value.String(), flag.Value.Type()))
// 						}
// 						valueNode.Content = append(valueNode.Content, seqNode)
// 					}
// 				}
// 				if valueNode.Kind == 0 {
// 					panic(fmt.Sprintf(
// 						"cannot marshal flag %q value %q type %q",
// 						flag.Name, flag.Value.String(), flag.Value.Type()))
// 				}
// 				curNode.Content = append(curNode.Content, keyNode, valueNode)
// 			}
// 		}
// 		prev = cmd
// 	}
// 	return yaml.NewEncoder(os.Stdout).Encode(topNode)
// }

// func (s *defaultsCommand) applyYamlScalarValue(node *yaml.Node, typ string, value string) {
// 	switch {
// 	case strings.HasPrefix(typ, "int") ||
// 		strings.HasPrefix(typ, "uint"):
// 		node.Kind = yaml.ScalarNode
// 		node.Tag = "!!int"
// 		node.Value = value
// 	case typ == "string" || typ == "duration":
// 		node.Kind = yaml.ScalarNode
// 		node.Tag = "!!str"
// 		node.Value = value
// 	case typ == "bool":
// 		node.Kind = yaml.ScalarNode
// 		node.Tag = "!!bool"
// 		node.Value = value
// 	}
// }

type MainCommand struct {
	*commands.BaseCommand
}

// func (s *MainCommand) InitCommand() error {
// 	s.Cobra().SilenceErrors = true
// 	s.Cobra().SilenceUsage = true
// 	s.Cobra().BashCompletionFunction = commands.CustomBashCompletionFunc(s.Cmd.Use)
// 	flags := &pflag.FlagSet{}
// 	flags.String("config", "", "Config file")
// 	flags.String("output", "human", "Output format (supported: json, yaml and human")
// 	flags.String("colours", "auto", "Use terminal colours (supported: auto, true, false)")
// 	flags.BoolP("debug", "d", false, "Enable log level debug")
// 	flags.Duration("client-timeout", 300*time.Second, "Timeout for requests")
// 	s.AddPersistentFlags(flags)

// 	configFile := config.AppConfig.GetString("config")

// 	fmt.Printf("0config file set %v\n", configFile)
// 	// Config load should be kept after main command flags declaration as --config
// 	// flag can override the config file path
// 	// if err := config.AppConfig.Load(); err != nil {
// 	// 	fmt.Fprintln(os.Stderr, err)
// 	// 	// return err
// 	// }

// 	if err := upapi.NewServiceClient(config.AppConfig); err != nil {
// 		return err
// 	}

// 	commands.BuildCommand(
// 		commands.New("completion", "Generate shell completion"),
// 		s, config.AppConfig,
// 	)

// 	commands.BuildCommand(
// 		&defaultsCommand{
// 			commands.New("defaults", "Generate defaults"),
// 		},
// 		s, config.AppConfig,
// 	)

// 	commands.BuildCommand(
// 		&versionCommand{commands.New("version", "Display software version")},
// 		s, config.AppConfig,
// 	)

// 	all.BuildCommands(s, s.Config())

// 	s.Cobra().SetUsageTemplate(ui.CommandUsageTemplate())
// 	s.Cobra().SetUsageFunc(ui.UsageFunc)

// 	return nil
// }

// func (s *MainCommand) MakePersistentPreExecuteCommand() func(args []string) error {
// 	return func(args []string) error {
// 		if err := validation.Value(s.Config().GetString("output"), "json", "yaml", "human"); err != nil {
// 			return fmt.Errorf("invalid output: %v", err)
// 		}
// 		if err := validation.Value(s.Config().GetString("colours"), "auto", "true", "false", "on", "off"); err != nil {
// 			return fmt.Errorf("invalid colours: %v", err)
// 		}
// 		switch strings.ToLower(s.Config().GetString("colours")) {
// 		case "on", "true":
// 			terminal.ForceColours(true)
// 		case "off", "false":
// 			terminal.ForceColours(false)
// 		}
// 		return nil
// 	}
// }

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", config.EnvPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func initializeConfig(cmd *cobra.Command, conf *config.Config, configFile string) error {
	v := conf.Viper

	v.SetConfigName(config.ConfigFileDefaultName)
	v.SetEnvPrefix(config.EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()
	v.SetConfigType("yaml")

	// Force config file in use if set
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
			return err
		}
		//Debug print CONFIG FILE NOT FOUND USING
	}

	// Update viper config name
	v.Set("config", v.ConfigFileUsed())

	// Bind main flags to viper
	bindFlags(cmd, v)

	if err := upapi.NewServiceClient(conf, v); err != nil {
		return err
	}

	return nil
}

// BuildRootCmd()
func BuildRootCmd(args []string, conf *config.Config) cobra.Command {
	var configPathFlag string
	var outputFlag string
	var colorsFlag string

	rootCmd := cobra.Command{
		Use:   "upctl",
		Short: "UpCloud CLI",
	}

	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.BashCompletionFunction = commands.CustomBashCompletionFunc(rootCmd.Use)

	flags := &pflag.FlagSet{}
	flags.StringVarP(&configPathFlag, "config", "c", "", "Config file")
	flags.StringVarP(&outputFlag, "output", "o", "human", "Output format (supported: json, yaml and human")
	flags.StringVar(&colorsFlag, "colours", "auto", "Use terminal colours (supported: auto, true, false)")
	// flags.BoolP("debug", "d", false, "Enable log level debug")
	// flags.Duration("client-timeout", 300*time.Second, "Timeout for requests")

	// Add flags
	flags.VisitAll(func(flag *pflag.Flag) {
		rootCmd.PersistentFlags().AddFlag(flag)
	})

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd, conf, configPathFlag)
	}

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier

	}
	_ = flags.Parse(args)

	// if err := initializeConfig(rootCmd, configPathFlag); err != nil {
	// 	return rootCmd, err
	// }

	// if err := upapi.NewServiceClient(conf, rootCmd.Viper); err != nil {
	// 	return rootCmd, err
	// }

	// commands.BuildCommand(
	// 	commands.New("completion", "Generate shell completion"),
	// 	s, config.AppConfig,
	// )

	// commands.BuildCommand(
	// 	&defaultsCommand{
	// 		commands.New("defaults", "Generate defaults"),
	// 	},
	// 	s, config.AppConfig,
	// )

	// commands.BuildCommand(
	// 	&versionCommand{commands.New("version", "Display software version")},
	// 	s, config.AppConfig,
	// )

	// all.BuildCommands(s, s.Config())

	// s.Cobra().SetUsageTemplate(ui.CommandUsageTemplate())
	// s.Cobra().SetUsageFunc(ui.UsageFunc)

	return rootCmd
}

//BootstrapCLI CLI entrypoint
func BootstrapCLI(args []string) error {

	conf := config.New()
	rootCmd := BuildRootCmd(args, conf)
	root.BuildAllCommands(conf, &rootCmd)

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
