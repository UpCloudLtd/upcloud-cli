package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func New(name, usage string) Command {
	return &baseCommand{
		Command: &cobra.Command{Use: name, Short: usage},
	}
}

type Command interface {
	SetConfig(viper *viper.Viper)
	SetParent(Command)
	Parent() Command
	Name() string
	InitCommand()
	MakeExecuteCommand() func(args []string) error
	MakePreExecuteCommand() func(args []string) error
	MakePersistentPreExecuteCommand() func(args []string) error
	HandleError(err error)
	HandleOutput(out interface{}) error
	CommandConfig
	CobraCommand
	Completion
}

type CommandConfig interface {
	Namespace() string
	Config() *viper.Viper
	SetConfigLoader(func(config *viper.Viper, loadContext int))
	ConfigLoader() func(config *viper.Viper, loadContext int)
	SetDefault(key string, val interface{})

	AddFlags(flag *pflag.FlagSet)
	AddVisibleColumnsFlag(flags *pflag.FlagSet, dstPtr *[]string, available, defaults []string)
	ConfigBindFlag(key string, flag *pflag.Flag)
	BoundFlags() map[string]*pflag.Flag

	ConfigValue(key string) interface{}
	ConfigValueString(key string) string
	FQConfigValue(key string) interface{}
	FQConfigValueString(key string) string
}

type Completion interface {
	ArgCompletion(fn func(toComplete string) ([]string, cobra.ShellCompDirective))
}

type CobraCommand interface {
	Cobra() *cobra.Command
}

const (
	ConfigLoadContextHelp = iota
	ConfigLoadContextRun
)

func BuildCommand(cmd, parent Command, config *viper.Viper) Command {
	cmd.SetParent(parent)
	cmd.SetConfig(config)
	cmd.InitCommand()
	if parent != nil {
		parent.Cobra().AddCommand(cmd.Cobra())
	}
	if cCmd := cmd.MakeExecuteCommand(); cCmd != nil && cmd.Cobra().RunE == nil {
		cmd.Cobra().RunE = func(_ *cobra.Command, args []string) error {
			if loader := cmd.ConfigLoader(); loader != nil {
				loader(config, ConfigLoadContextRun)
			}
			return cCmd(args)
		}
	}
	if cCmd := cmd.MakePreExecuteCommand(); cCmd != nil && cmd.Cobra().PreRunE == nil {
		cmd.Cobra().PreRunE = func(_ *cobra.Command, args []string) error {
			if loader := cmd.ConfigLoader(); loader != nil {
				loader(config, ConfigLoadContextRun)
			}
			return cCmd(args)
		}
	}
	if cCmd := cmd.MakePersistentPreExecuteCommand(); cCmd != nil && cmd.Cobra().PersistentPreRunE == nil {
		cmd.Cobra().PersistentPreRunE = func(_ *cobra.Command, args []string) error {
			if loader := cmd.ConfigLoader(); loader != nil {
				loader(config, ConfigLoadContextRun)
			}
			return cCmd(args)
		}
	}
	if len(cmd.BoundFlags()) > 0 {
		curHelp := cmd.Cobra().HelpFunc()
		cmd.Cobra().SetHelpFunc(func(cCmd *cobra.Command, args []string) {
			if loader := cmd.ConfigLoader(); loader != nil {
				loader(config, ConfigLoadContextHelp)
			}
			for k, v := range cmd.BoundFlags() {
				if !config.IsSet(k) {
					continue
				}
				v.DefValue = config.GetString(k)
			}
			curHelp(cCmd, args)
		})
	}
	return cmd
}

type baseCommand struct {
	*cobra.Command
	parent       Command
	config       *viper.Viper
	configLoader func(config *viper.Viper, loadContext int)
	boundFlags   map[string]*pflag.Flag
}

func (s *baseCommand) Name() string {
	return s.Command.Use
}

func (s *baseCommand) SetConfig(config *viper.Viper) {
	s.config = config
}

func (s *baseCommand) SetParent(command Command) {
	s.parent = command
}

func (s *baseCommand) Parent() Command {
	return s.parent
}

// A hook to handle flag registration.
// The config values are not available during this hook. Register a cobra hook to use them. You can set defaults though.
func (s *baseCommand) InitCommand() {
}

func (s *baseCommand) MakeExecuteCommand() func(args []string) error {
	return nil
}

func (s *baseCommand) MakePreExecuteCommand() func(args []string) error {
	return nil
}

func (s *baseCommand) MakePersistentPreExecuteCommand() func(args []string) error {
	return nil
}

// Returns the namespace of this command from the chain of parent commands
// The format is cmdRoot.child1.child2.childN
// No namespace is returned for the root command (parent == nil)
func (s *baseCommand) Namespace() string {
	var (
		sb    strings.Builder
		names []string
	)
	for c := s.parent; c != nil; c = c.Parent() {
		names = append(names, c.Name())
	}
	for i := len(names) - 1; i >= 0; i-- {
		sb.WriteString(names[i])
		sb.WriteString(".")
	}
	if s.parent != nil {
		sb.WriteString(s.Name())
	}
	return sb.String()
}

func (s *baseCommand) Cobra() *cobra.Command {
	return s.Command
}

// Config //

// Sets a default value for a config key
// The command namespace is appended before the key
func (s *baseCommand) SetDefault(key string, val interface{}) {
	s.config.SetDefault(s.Namespace()+"."+key, val)
}

func (s *baseCommand) ConfigValue(key string) interface{} {
	if s.Namespace() != "" {
		key = s.Namespace() + "." + key
	}
	return s.config.Get(key)
}

func (s *baseCommand) ConfigValueString(key string) string {
	if s.Namespace() != "" {
		key = s.Namespace() + "." + key
	}
	return s.config.GetString(key)
}

func (s *baseCommand) FQConfigValue(key string) interface{} {
	return s.config.Get(key)
}

func (s *baseCommand) FQConfigValueString(key string) string {
	return s.config.GetString(key)
}

func (s *baseCommand) Config() *viper.Viper {
	return s.config
}

// Adds a flag to the command and binds config value into it with namespace
func (s *baseCommand) AddFlags(flags *pflag.FlagSet) {
	if flags == nil {
		panic("Nil flagset")
	}
	flags.VisitAll(func(flag *pflag.Flag) {
		s.Cobra().Flags().AddFlag(flag)
		s.ConfigBindFlag(flag.Name, flag)
	})
}

func (s *baseCommand) AddVisibleColumnsFlag(flags *pflag.FlagSet, dstPtr *[]string, available, defaults []string) {
	flags.StringSliceVarP(dstPtr, "columns", "c", nil,
		fmt.Sprintf("Reorder or show additional columns in human readable output\navailable: %s",
			strings.Join(available, ",")))
	curPreRun := s.Cobra().PreRunE
	s.Cobra().PreRunE = func(cmd *cobra.Command, args []string) error {
		if curPreRun != nil {
			if err := curPreRun(cmd, args); err != nil {
				return err
			}
		}
		if !cmd.Flags().Changed("columns") {
			*dstPtr = defaults
		}

		return nil
	}
}

func (s *baseCommand) BoundFlags() map[string]*pflag.Flag {
	return s.boundFlags
}

func (s *baseCommand) ConfigBindFlag(key string, flag *pflag.Flag) {
	if flag == nil {
		panic("Nil flag bound")
	}
	if s.boundFlags == nil {
		s.boundFlags = make(map[string]*pflag.Flag)
	}
	s.boundFlags[key] = flag
	if s.Namespace() != "" {
		key = s.Namespace() + "." + key
	}
	_ = s.config.BindPFlag(key, flag)
}

func (s *baseCommand) SetConfigLoader(fn func(config *viper.Viper, loadContext int)) {
	s.configLoader = fn
}

func (s *baseCommand) ConfigLoader() func(config *viper.Viper, loadContext int) {
	return s.configLoader
}

// Error handling //
func (s *baseCommand) HandleError(err error) {
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())
	switch s.Config().GetString("output") {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		if isTerminal {
			enc.SetIndent("", "    ")
		}
		if ucApiErr, ok := err.(*upcloud.Error); ok {
			enc.Encode(ucApiErr)
			break
		}
		_ = enc.Encode(map[string]interface{}{"error": fmt.Sprintf("%v", err)})
	case "yaml":
		if ucApiErr, ok := err.(*upcloud.Error); ok {
			tmpMap := make(map[string]interface{})
			if b, err := json.Marshal(ucApiErr); err == nil {
				if err := json.Unmarshal(b, &tmpMap); err == nil {
					yaml.NewEncoder(os.Stdout).Encode(tmpMap)
					break
				}
			}
		}
		_ = yaml.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": fmt.Sprintf("%v", err)})
	default:
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	}
}

// Output handling //
func (s *baseCommand) HandleOutput(out interface{}) error {
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())
	switch s.Config().GetString("output") {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		if isTerminal {
			enc.SetIndent("", "    ")
		}
		enc.Encode(out)
	case "yaml":
		yaml.NewEncoder(os.Stdout).Encode(out)
	default:
		fmt.Printf("%v", out)
	}
	return nil
}

// Completion //
func (s *baseCommand) ArgCompletion(fn func(toComplete string) ([]string, cobra.ShellCompDirective)) {
	s.Cobra().ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return fn(toComplete)
	}
}

func MatchStringPrefix(vals []string, toComplete string) []string {
	var r []string
	if toComplete == "" {
		return vals
	}
	for _, v := range vals {
		if strings.HasPrefix(v, toComplete) {
			r = append(r, v)
		}
	}
	return r
}
