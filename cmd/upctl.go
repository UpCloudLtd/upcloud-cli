package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/all"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/terminal"
	"github.com/UpCloudLtd/cli/internal/ui"
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

func (s *completionCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		switch args[0] {
		case "bash":
			_ = s.Cobra().Root().GenBashCompletion(os.Stdout)
		}
		return nil, nil
	}
}

type defaultsCommand struct {
	*commands.BaseCommand
}

func (s *defaultsCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		output, err := s.HandleOutput(nil)
		if err != nil {
			return nil, err
		}
		fmt.Println(output)
		return nil, nil
	}
}

func (s *defaultsCommand) HandleOutput(_ interface{}) (string, error) {
	output := config.ConfigValueOutputYaml
	if s.Config().Top().IsSet(config.ConfigKeyOutput) && !s.Config().OutputHuman() {
		output = s.Config().Output()
	}
	if output == config.ConfigValueOutputJson {
		return "", fmt.Errorf("only yaml output is supported for this command")
	}

	var commandsWithFlags []commands.Command
nextChild:
	for _, cmd := range append([]commands.Command{s.Parent()}, s.Parent().Children()...) {
		if len(cmd.Config().BoundFlags()) > 0 {
			commandsWithFlags = append(commandsWithFlags, cmd)
		} else {
			for _, ccmd := range cmd.Children() {
				if len(ccmd.Config().BoundFlags()) > 0 {
					commandsWithFlags = append(commandsWithFlags, cmd)
					continue nextChild
				}
			}
		}
	}
	topNode := &yaml.Node{Kind: yaml.DocumentNode}
	topNode.HeadComment = "UpCtl configuration defaults"
	curNode := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	topNode.Content = append(topNode.Content, curNode)
	commandToNode := map[commands.Command]*yaml.Node{s.Parent(): curNode}
	var prev commands.Command
	for _, cmd := range commandsWithFlags {
		if prev != nil && prev != cmd.Parent() {
			curNode = commandToNode[cmd.Parent()]
		}
		// Skip top level command name
		if cmd.Parent() != nil {
			curNode.Content = append(curNode.Content,
				&yaml.Node{
					Kind:  yaml.ScalarNode,
					Tag:   "!!str",
					Value: cmd.Name(),
				})
			curNode.Content = append(curNode.Content, &yaml.Node{
				Kind: yaml.MappingNode,
				Tag:  "!!map",
			})
			curNode = curNode.Content[len(curNode.Content)-1]
		}
		commandToNode[cmd] = curNode
		if len(cmd.Config().BoundFlags()) > 0 {
			for _, flag := range cmd.Config().BoundFlags() {
				keyNode := &yaml.Node{
					Kind:        yaml.ScalarNode,
					Tag:         "!!str",
					HeadComment: flag.Usage,
					Value:       flag.Name,
					FootComment: "\r",
				}
				valueNode := &yaml.Node{}
				s.applyYamlScalarValue(valueNode, flag.Value.Type(), flag.DefValue)
				if valueNode.Kind == 0 && strings.HasSuffix(flag.Value.Type(), "Slice") {
					styp := flag.Value.Type()[0 : len(flag.Value.Type())-5]
					valueNode.Kind = yaml.SequenceNode
					valueNode.Style = yaml.FlowStyle
					valueNode.Tag = "!!seq"
					flagSvals := strings.TrimPrefix(flag.Value.String(), "[")
					flagSvals = strings.TrimSuffix(flagSvals, "]")
					for _, fsv := range strings.Split(flagSvals, ",") {
						seqNode := &yaml.Node{}
						s.applyYamlScalarValue(seqNode, styp, fsv)
						if seqNode.Kind == 0 {
							panic(fmt.Sprintf(
								"cannot marshal %q value %q type %q",
								flag.Name, flag.Value.String(), flag.Value.Type()))
						}
						valueNode.Content = append(valueNode.Content, seqNode)
					}
				}
				if valueNode.Kind == 0 {
					panic(fmt.Sprintf(
						"cannot marshal flag %q value %q type %q",
						flag.Name, flag.Value.String(), flag.Value.Type()))
				}
				curNode.Content = append(curNode.Content, keyNode, valueNode)
			}
		}
		prev = cmd
	}
	buf := new(bytes.Buffer)
	err := yaml.NewEncoder(buf).Encode(topNode)
	if err != nil {
		return "", nil
	}
	return buf.String(), nil
}

func (s *defaultsCommand) applyYamlScalarValue(node *yaml.Node, typ string, value string) {
	switch {
	case strings.HasPrefix(typ, "int") ||
		strings.HasPrefix(typ, "uint"):
		node.Kind = yaml.ScalarNode
		node.Tag = "!!int"
		node.Value = value
	case typ == "string" || typ == "duration":
		node.Kind = yaml.ScalarNode
		node.Tag = "!!str"
		node.Value = value
	case typ == "bool":
		node.Kind = yaml.ScalarNode
		node.Tag = "!!bool"
		node.Value = value
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
	s.Config().Viper().SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
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
	s.Cobra().BashCompletionFunction = commands.CustomBashCompletionFunc(s.Name())
	flags := &pflag.FlagSet{}
	flags.String("config", "", "Config file")
	flags.String("output", "human", "Output format (supported: json, yaml and human")
	flags.String("colours", "auto", "Use terminal colours (supported: auto, true, false)")
	flags.Duration("client-timeout", 600*time.Second, "Timeout for requests")
	s.AddPersistentFlags(flags)

	s.SetConfigLoader(func(config *config.Config, loadContext int) {
		s.initConfig(loadContext == commands.ConfigLoadContextHelp)
	})

	commands.BuildCommand(
		&completionCommand{commands.New("completion", "Generate shell completion code")},
		s, config.New(mainConfig.Viper()))

	commands.BuildCommand(
		&defaultsCommand{commands.New("defaults", "Generate defaults")},
		s, config.New(mainConfig.Viper()))

	if loader := s.ConfigLoader(); loader != nil {
		loader(s.Config(), commands.ConfigLoadContextRun)
	}

	all.BuildCommands(s, s.Config())

	s.Cobra().SetUsageTemplate(ui.CommandUsage())
	s.Cobra().SetUsageFunc(ui.UsageFunc)
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
