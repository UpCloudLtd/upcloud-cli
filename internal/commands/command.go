package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"

	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/ui"
)

func New(name, usage string) *BaseCommand {
	return &BaseCommand{
		name:    name,
		Command: &cobra.Command{Use: name, Short: usage},
	}
}

type Command interface {
	SetConfig(config *config.Config)
	SetParent(Command)
	Parent() Command
	Name() string
	InitCommand()
	MakeExecuteCommand() func(args []string) error
	MakePreExecuteCommand() func(args []string) error
	MakePersistentPreExecuteCommand() func(args []string) error
	SetConfigLoader(func(config *config.Config, loadContext int))
	ConfigLoader() func(config *config.Config, loadContext int)
	Config() *config.Config
	HandleOutput(out interface{}) error
	HandleError(err error)
	CobraCommand
}

type CobraCommand interface {
	Cobra() *cobra.Command
}

type namespace interface {
	Namespace() string
}

const (
	ConfigLoadContextHelp = iota
	ConfigLoadContextRun
)

func BuildCommand(child, parent Command, config *config.Config) Command {
	child.SetParent(parent)
	child.SetConfig(config)
	if parent != nil {
		child.SetConfigLoader(parent.ConfigLoader())
	}
	if nsCmd, ok := child.(namespace); ok {
		config.SetNamespace(nsCmd.Namespace())
	}
	child.InitCommand()
	if cCmd := child.MakeExecuteCommand(); cCmd != nil && child.Cobra().RunE == nil {
		child.Cobra().RunE = func(_ *cobra.Command, args []string) error {
			if loader := child.ConfigLoader(); loader != nil {
				loader(config, ConfigLoadContextRun)
			}
			return cCmd(args)
		}
	}
	if cCmd := child.MakePreExecuteCommand(); cCmd != nil && child.Cobra().PreRunE == nil {
		child.Cobra().PreRunE = func(_ *cobra.Command, args []string) error {
			if loader := child.ConfigLoader(); loader != nil {
				loader(config, ConfigLoadContextRun)
			}
			return cCmd(args)
		}
	}
	if cCmd := child.MakePersistentPreExecuteCommand(); cCmd != nil && child.Cobra().PersistentPreRunE == nil {
		child.Cobra().PersistentPreRunE = func(_ *cobra.Command, args []string) error {
			if loader := child.ConfigLoader(); loader != nil {
				loader(config, ConfigLoadContextRun)
			}
			return cCmd(args)
		}
	}

	curHelp := child.Cobra().HelpFunc()
	child.Cobra().SetHelpFunc(func(cCmd *cobra.Command, args []string) {
		if loader := child.ConfigLoader(); loader != nil {
			loader(config, ConfigLoadContextHelp)
		}
		for cmd := child; cmd != nil; cmd = cmd.Parent() {
			for k, v := range cmd.Config().BoundFlags() {
				if !cmd.Config().Viper().IsSet(k) {
					continue
				}
				v.DefValue = cmd.Config().GetString(k)
			}
		}
		curHelp(cCmd, args)
	})

	// Need to set child command in the end as otherwise HelpFunc() returns the parent's helpfunc
	if parent != nil {
		parent.Cobra().AddCommand(child.Cobra())
	}
	return child
}

type BaseCommand struct {
	*cobra.Command
	name         string
	parent       Command
	config       *config.Config
	configLoader func(config *config.Config, loadContext int)
}

func (s *BaseCommand) Name() string {
	return s.name
}

func (s *BaseCommand) SetConfig(config *config.Config) {
	s.config = config
}

func (s *BaseCommand) SetParent(command Command) {
	s.parent = command
}

func (s *BaseCommand) Parent() Command {
	return s.parent
}

// A hook to handle flag registration.
// The config values are not available during this hook. Register a cobra hook to use them. You can set defaults though.
func (s *BaseCommand) InitCommand() {
}

func (s *BaseCommand) MakeExecuteCommand() func(args []string) error {
	return nil
}

func (s *BaseCommand) MakePreExecuteCommand() func(args []string) error {
	return nil
}

func (s *BaseCommand) MakePersistentPreExecuteCommand() func(args []string) error {
	return nil
}

// Returns the namespace of this command from the chain of parent commands
// The format is cmdRoot.child1.child2.childN
// No namespace is returned for the root command (parent == nil)
func (s *BaseCommand) Namespace() string {
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

func (s *BaseCommand) Cobra() *cobra.Command {
	return s.Command
}

// Config //

func (s *BaseCommand) Config() *config.Config {
	return s.config
}

func (s *BaseCommand) SetConfigLoader(fn func(config *config.Config, loadContext int)) {
	s.configLoader = fn
}

func (s *BaseCommand) ConfigLoader() func(config *config.Config, loadContext int) {
	return s.configLoader
}

// Flags //

// Adds a flagset to the command and binds config value into it with namespace
func (s *BaseCommand) AddFlags(flags *pflag.FlagSet) {
	if flags == nil {
		panic("Nil flagset")
	}
	flags.VisitAll(func(flag *pflag.Flag) {
		flag.Usage = text.WrapSoft(flag.Usage, ui.CommandUsageLineLength)
		s.Cobra().Flags().AddFlag(flag)
		s.config.ConfigBindFlag(flag.Name, flag)
	})
}

// Adds a persistent flagset to the command and binds config value into it with namespace
func (s *BaseCommand) AddPersistentFlags(flags *pflag.FlagSet) {
	if flags == nil {
		panic("Nil flagset")
	}
	flags.VisitAll(func(flag *pflag.Flag) {
		flag.Usage = text.WrapSoft(flag.Usage, ui.CommandUsageLineLength)
		s.Cobra().PersistentFlags().AddFlag(flag)
		s.config.ConfigBindFlag(flag.Name, flag)
	})
}

func (s *BaseCommand) AddVisibleColumnsFlag(flags *pflag.FlagSet, dstPtr *[]string, available, defaults []string) {
	flags.StringSliceVarP(dstPtr, "columns", "c", nil,
		text.WrapSoft(fmt.Sprintf("Reorder or show additional columns in human readable output\navailable: %s",
			strings.Join(available, ",")), ui.CommandUsageLineLength))
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

func (s *BaseCommand) SetPositionalArgHelp(help string) {
	if help == "" {
		s.Use = s.name
		return
	}
	s.Use = fmt.Sprintf("%s %s", s.name, help)
}

// Error handling //
func (s *BaseCommand) HandleError(err error) {
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())
	switch s.Config().GetString("output") {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		if isTerminal {
			enc.SetIndent("", "    ")
		}
		if ucApiErr, ok := err.(*upcloud.Error); ok {
			_ = enc.Encode(ucApiErr)
			break
		}
		_ = enc.Encode(map[string]interface{}{"error": fmt.Sprintf("%v", err)})
	case "yaml":
		if ucApiErr, ok := err.(*upcloud.Error); ok {
			tmpMap := make(map[string]interface{})
			if b, err := json.Marshal(ucApiErr); err == nil {
				if err := json.Unmarshal(b, &tmpMap); err == nil {
					_ = yaml.NewEncoder(os.Stdout).Encode(tmpMap)
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
func (s *BaseCommand) HandleOutput(out interface{}) error {
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())
	switch s.Config().Output() {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		if isTerminal {
			enc.SetIndent("", "    ")
		}
		_ = enc.Encode(out)
	case "yaml":
		// TODO(aakso): maybe we need to patch the yaml library to get field names from json tags?
		//              that will doubtly get accepted though.
		_ = yaml.NewEncoder(os.Stdout).Encode(out)
	default:
		fmt.Printf("%v", out)
	}
	return nil
}

// Completion //
func (s *BaseCommand) ArgCompletion(fn func(toComplete string) ([]string, cobra.ShellCompDirective)) {
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
