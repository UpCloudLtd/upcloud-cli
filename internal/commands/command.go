package commands

import (
	"encoding/json"
	"fmt"
	// "io"
	"os"
	// "sort"
	"strings"

	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	gyaml "github.com/ghodss/yaml"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func GenerateCmdList(cmds ...*BaseCommand) []*BaseCommand {
	var c = []*BaseCommand{}
	for _, cmd := range cmds {
		c = append(c, cmd)
	}

	return c
}

// New returns a BaseCommand that implements Command. It is used as a base to create custom commands from.
func New(name, usage string) *BaseCommand {
	return &BaseCommand{
		Cmd: &cobra.Command{Use: name, Short: usage},
	}
}

// Command defines the common functionality for all commands
// BaseCommand is the base type for all commands, implementing Command
type BaseCommand struct {
	// Name is the command name
	Name string

	// Aliases is an array of aliases that can be used instead of the first word in Use.
	Aliases []string

	// Short is the short description shown in the 'help' output.
	Short string

	// Long is the long message shown in the 'help <this-command>' output.
	Long string

	Cmd              *cobra.Command
	Viper            *viper.Viper
	Parent           *BaseCommand
	Cfg              *config.Config
	Run              RunCmd
	PreRun           PreRunCmd
	PersistentPreRun PersistentPreRunCmd
	// HandleOutput     HandleOutputCmd
	// HandleError      HandleErrorCmd
}
type RunCmd func(conf *config.Config, args interface{}) (i interface{}, e error)
type PreRunCmd func(args []string) error
type PersistentPreRunCmd func(args []string) error

// type HandleOutputCmd func(writer io.Writer, out interface{}) error
type HandleErrorCmd func(err error)

func runCmd(conf *config.Config, cmd *BaseCommand) func(*cobra.Command, []string) error {
	return func(cobraCmd *cobra.Command, args []string) error {
		if cmd.Run == nil {
			return nil
		}

		res, err := cmd.Run(conf, args)
		if err != nil {
			return err
		}

		return handleOutput(res, conf.Output())
	}
}

// BuildCommand sets up a Command with the specified config and adds it to Cobra
func BuildCommand(conf *config.Config, cmd *BaseCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   cmd.Name,
		Short: cmd.Short,
		Long:  cmd.Long,
	}

	// Config
	cobraCmd.Flags().SortFlags = false

	cobraCmd.RunE = runCmd(conf, cmd)

	// // Apply values set from viper

	// child.Cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
	// 	for cmd := child; cmd != nil; cmd = cmd.Parent {
	// 		for _, v := range cmd.Config().BoundFlags() {
	// 			if !cmd.Config().IsSet(v.Name) {
	// 				continue
	// 			}
	// 			if v.Changed {
	// 				continue
	// 			}
	// 			if err := v.Value.Set(cmd.Config().GetString(v.Name)); err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// 	return nil
	// }

	// if cCmd := child.MakePreExecuteCommand(); cCmd != nil && child.Cobra().PreRunE == nil {
	// 	child.Cobra().PreRunE = func(_ *cobra.Command, args []string) error {
	// 		if loader := child.ConfigLoader(); loader != nil {
	// 			if err := loader(config); err != nil {
	// 				return fmt.Errorf("Config load: %v", err)
	// 			}
	// 		}
	// 		return cCmd(args)
	// 	}
	// }
	// if cCmd := child.MakePersistentPreExecuteCommand(); cCmd != nil && child.Cobra().PersistentPreRunE == nil {
	// 	child.Cobra().PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
	// 		if loader := child.ConfigLoader(); loader != nil {
	// 			if err := loader(config); err != nil {
	// 				return fmt.Errorf("Config load: %v", err)
	// 			}
	// 		}
	// 		return cCmd(args)
	// 	}
	// }

	// curHelp := child.Cobra().HelpFunc()
	// child.Cmd.SetHelpFunc(func(cCmd *cobra.Command, args []string) {
	// 	for cmd := child; cmd != nil; cmd = cmd.Parent {
	// 		for _, v := range cmd.Config().BoundFlags() {
	// 			if !cmd.Config().IsSet(v.Name) {
	// 				continue
	// 			}
	// 			v.DefValue = cmd.Config().GetString(v.Name)
	// 		}
	// 	}
	// 	curHelp(cCmd, args)
	// })

	return cobraCmd
}

// BaseCommand is the base type for all commands, implementing Command
// type BaseCommand struct {
// 	Cmd              *cobra.Command
// 	name             string
// 	parent           Command
// 	childrenPos      map[Command]int
// 	nextChildSortPos int
// 	Cfg              *config.Config
// 	configLoader     func(config *config.Config) error
// }

// // Name returns the name of the command
// func (s *BaseCommand) Name() string {
// 	return s.name
// }

// // SetChild sets command as the child of this command
// func (s *BaseCommand) SetChild(command Command) {
// 	if command == nil {
// 		return
// 	}
// 	if _, alreadyChild := s.childrenPos[command]; alreadyChild {
// 		return
// 	}
// 	if s.childrenPos == nil {
// 		s.childrenPos = make(map[Command]int)
// 	}
// 	s.childrenPos[command] = s.nextChildSortPos
// 	s.nextChildSortPos++
// 	if command.Parent() != s {
// 		command.SetParent(s)
// 	}
// }

// // DeleteChild removes command from the children of this command
// func (s *BaseCommand) DeleteChild(command Command) {
// 	if command.Parent() == s {
// 		command.SetParent(nil)
// 	}
// 	delete(s.childrenPos, command)
// }

// // Children returns a list of all the child commands of this command (eg. including the children of children)
// func (s *BaseCommand) Children() []Command {
// 	var (
// 		r      []Command
// 		sorted []Command
// 	)
// 	for child := range s.childrenPos {
// 		sorted = append(sorted, child)
// 	}
// 	sort.Slice(sorted, func(i, j int) bool {
// 		return s.childrenPos[sorted[i]] < s.childrenPos[sorted[j]]
// 	})
// 	for _, child := range sorted {
// 		r = append(r, child)
// 		r = append(r, child.Children()...)
// 	}
// 	return r
// }

// // SetParent sets the parent of this command to given command
// func (s *BaseCommand) SetParent(command Command) {
// 	if s.parent != nil {
// 		s.DeleteChild(command)
// 	}
// 	s.parent = command
// 	if s.parent != nil {
// 		s.parent.SetChild(s)
// 	}
// }

// // Parent returns the parent of the command
// func (s *BaseCommand) Parent() Command {
// 	return s.parent
// }

// SetFlags parses the given flags
func (s *BaseCommand) SetFlags(flags []string) error {
	return s.Cobra().Flags().Parse(flags)
}

// InitCommand can be overriden to handle flag registration.
// A hook to handle flag registration.
// The config values are not available during this hook. Register a cobra hook to use them. You can set defaults though.
func (s *BaseCommand) InitCommand() error {
	return nil
}

// MakeExecuteCommand should be overwritten by the actual command implementations
// The function returned is ran during the 'regular' execute phase and the returned value is
// returned to the user, formatted as requested.
func (s *BaseCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return nil
}

// MakePreExecuteCommand should be overwritten by the actual command implementations
// The function returned is ran before the 'regular' execute phase.
func (s *BaseCommand) MakePreExecuteCommand() func(args []string) error {
	return nil
}

// MakePersistentPreExecuteCommand should be overwritten by the actual command implementations
// The function returned is ran before PreExecuteCommand().
func (s *BaseCommand) MakePersistentPreExecuteCommand() func(args []string) error {
	return nil
}

// // Namespace returns the namespace of this command from the chain of parent commands
// // The format is cmdRoot.child1.child2.childN
// // No namespace is returned for the root command (parent == nil)
// func (s *BaseCommand) Namespace() string {
// 	var (
// 		sb    strings.Builder
// 		names []string
// 	)
// 	for c := s.parent; c != nil; c = c.Parent() {
// 		// Skip root command name in namespace
// 		if c.Parent() == nil {
// 			continue
// 		}
// 		names = append(names, c.Name())
// 	}
// 	for i := len(names) - 1; i >= 0; i-- {
// 		sb.WriteString(names[i])
// 		sb.WriteString(".")
// 	}
// 	if s.parent != nil {
// 		sb.WriteString(s.Name())
// 	}
// 	return sb.String()
// }

// Cobra returns the underlying *cobra.Command
func (s *BaseCommand) Cobra() *cobra.Command {
	return s.Cmd
}

// Config //

// Config implements Command.Config, returning the *config.Config of the command
func (s *BaseCommand) Config() *config.Config {
	return s.Cfg
}

// // SetConfigLoader implements Command.SetConfigLoader, setting internal config loader
// func (s *BaseCommand) SetConfigLoader(fn func(config *config.Config) error) {
// 	s.configLoader = fn
// }

// // ConfigLoader implements Command.ConfigLoader, returning the specified ConfigLoader
// func (s *BaseCommand) ConfigLoader() func(config *config.Config) error {
// 	return s.configLoader
// }

// Flags //

// // AddFlags adds a flagset to the command and binds config value into it with namespace
// func (s *BaseCommand) AddFlags(flags *pflag.FlagSet) {
// 	if flags == nil {
// 		panic("Nil flagset")
// 	}
// 	flags.VisitAll(func(flag *pflag.Flag) {
// 		s.Cobra().Flags().AddFlag(flag)
// 	})
// 	s.Cfg.ConfigBindFlagSet(flags)
// }

// // AddPersistentFlags adds a persistent flagset to the command and binds config value into it with namespace
// func (s *BaseCommand) AddPersistentFlags(flags *pflag.FlagSet) {
// 	if flags == nil {
// 		panic("Nil flagset")
// 	}
// 	flags.VisitAll(func(flag *pflag.Flag) {
// 		s.Cmd.PersistentFlags().AddFlag(flag)
// 	})
// 	s.Cfg.ConfigBindFlagSet(flags)
// }

// AddVisibleColumnsFlag is a convenience method to set a common flag '--columns' to commands
func (s *BaseCommand) AddVisibleColumnsFlag(flags *pflag.FlagSet, dstPtr *[]string, available, defaults []string) {
	flags.StringSliceVarP(dstPtr, "columns", "c", defaults,
		fmt.Sprintf("Reorder or show additional columns in human readable output.\nAvailable: %s",
			strings.Join(available, ",")))
}

// SetPositionalArgHelp is a convenience method to set the help text for positional arguments.
// if help is an empty string, uses just the name of the command.
func (s *BaseCommand) SetPositionalArgHelp(help string) {
	if help == "" {
		return
	}
	s.Cmd.Use = fmt.Sprintf("%s %s", s.Cmd.Use, help)
}

// Error handling //

// HandleError is used to handle errors from the main command
func (s *BaseCommand) HandleError(err error) {
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())
	switch s.Config().GetString("output") {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		if isTerminal {
			enc.SetIndent("", "    ")
		}
		if ucAPIErr, ok := err.(*upcloud.Error); ok {
			_ = enc.Encode(ucAPIErr)
			break
		}
		_ = enc.Encode(map[string]interface{}{"error": fmt.Sprintf("%v", err)})
	case "yaml":
		if ucAPIErr, ok := err.(*upcloud.Error); ok {
			tmpMap := make(map[string]interface{})
			if b, err := json.Marshal(ucAPIErr); err == nil {
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

func handleOutput(out interface{}, format string) error {
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())
	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		if isTerminal {
			enc.SetIndent("", "    ")
		}
		_ = enc.Encode(out)
	case "yaml":
		bytes, err := gyaml.Marshal(out)
		if err != nil {
			return err
		}
		_, _ = os.Stdout.Write(bytes)
	default:
		fmt.Printf("%v", out)
	}
	return nil
}

// Completion //

// ArgCompletion is a convenience method to set upctl-specific completion function to the underlying cobra.Command
func (s *BaseCommand) ArgCompletion(fn func(toComplete string) ([]string, cobra.ShellCompDirective)) {
	s.Cobra().ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return fn(toComplete)
	}
}

// MatchStringPrefix returns a list of string in vals which have a prefix as specified in key. Quotes are removed from key and output strings are escaped according to completion rules
func MatchStringPrefix(vals []string, key string, caseSensitive bool) []string {
	var r []string
	key = strings.TrimPrefix(key, `"`)
	key = strings.TrimPrefix(key, "'")
	key = strings.TrimSuffix(key, `"`)
	key = strings.TrimSuffix(key, "'")
	for _, v := range vals {
		if (caseSensitive && strings.HasPrefix(v, key)) ||
			(!caseSensitive && strings.HasPrefix(strings.ToLower(v), strings.ToLower(key))) ||
			key == "" {
			r = append(r, CompletionEscape(v))
		}
	}
	return r
}

// CompletionEscape escapes a string according to completion rules (?)
// in effect, this means that the string will be quoted with double quotes if it contains a space or parentheses.
func CompletionEscape(s string) string {
	if strings.ContainsAny(s, ` ()`) {
		return fmt.Sprintf(`"%s"`, s)
	}
	return s
}
