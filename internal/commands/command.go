package commands

import (
	"encoding/json"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/mapper"
	"github.com/UpCloudLtd/cli/internal/output"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	gyaml "github.com/ghodss/yaml"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

// New returns a BaseCommand that implements Command. It is used as a base to create custom commands from.
func New(name, usage string) *BaseCommand {
	return &BaseCommand{
		name:  name,
		cobra: &cobra.Command{Use: name, Short: usage},
	}
}

// Command defines the common functionality for all commands
type Command interface {
	SetConfig(config *config.Config)
	SetParent(Command)
	SetChild(command Command)
	SetFlags(flags []string) error
	DeleteChild(command Command)
	Children() []Command
	Parent() Command
	Name() string
	InitCommand()
	MakeExecuteCommand() func(args []string) (interface{}, error)
	MakePreExecuteCommand() func(args []string) error
	MakePersistentPreExecuteCommand() func(args []string) error
	SetConfigLoader(func(config *config.Config) error)
	ConfigLoader() func(config *config.Config) error
	Config() *config.Config
	HandleOutput(writer io.Writer, out interface{}) error
	HandleError(err error)
	CobraCommand
}

// NewCommand is a new container for commands, currently still including the old interface until we can deprecate it
type NewCommand interface {
	Execute(exec Executor, arg string) (output.Command, error)
	MaximumExecutions() int
	ArgumentMapper() (mapper.Argument, error)
	NewParent() NewCommand
	Command
}

// CobraCommand is an interface for commands that can refer back to their base cobra.Command
type CobraCommand interface {
	Cobra() *cobra.Command
}

type namespace interface {
	Namespace() string
}

// BuildCommand sets up a Command with the specified config and adds it to Cobra
func BuildCommand(child, parent Command, config *config.Config) Command {
	child.SetParent(parent)
	child.SetConfig(config)
	child.Cobra().Flags().SortFlags = false
	if parent != nil {
		child.SetConfigLoader(parent.ConfigLoader())
	}
	if nsCmd, ok := child.(namespace); ok {
		config.SetNamespace(nsCmd.Namespace())
	}
	child.InitCommand()
	// Apply values set from viper
	child.Cobra().PreRunE = func(cmd *cobra.Command, args []string) error {
		for cmd := child; cmd != nil; cmd = cmd.Parent() {
			for _, v := range cmd.Config().BoundFlags() {
				if !cmd.Config().IsSet(v.Name) {
					continue
				}
				if v.Changed {
					continue
				}
				if err := v.Value.Set(cmd.Config().GetString(v.Name)); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if nc, ok := child.(NewCommand); ok {
		child.Cobra().RunE = func(cmd *cobra.Command, args []string) error {
			executor := NewExecutor(config)
			mapper, err := nc.ArgumentMapper()
			if err != nil {
				return fmt.Errorf("invalid mapper: %w", err)
			}
			returnChan := make(chan executeResult)
			for i, arg := range args {
				go func(index int, argument string) {
					executeArgument := argument
					if mapper != nil {
						if res, err := mapper(argument); err == nil {
							executeArgument = res
						} else {
							executor.NewLogEntry(fmt.Sprintf("invalid argument: %v", err))
							returnChan <- executeResult{Job: index, Error: fmt.Errorf("cannot map argument '%v': %w", argument, err)}
							return
						}
					}
					res, err := nc.Execute(executor, executeArgument)
					returnChan <- executeResult{Job: index, Result: res, Error: err}
				}(i, arg)
			}
			result := map[int]executeResult{}
			renderTicker := time.NewTicker(100 * time.Millisecond)
		waitForResults:
			for {
				select {
				case res := <-returnChan:
					result[res.Job] = res
					if len(result) == len(args) {
						// we're done
						break waitForResults
					}
				case <-renderTicker.C:
					if config.InteractiveUI() {
						executor.Update()
					}
				}
			}
			executor.Update()
			if len(result) > 1 {
				resultList := []interface{}{}
				for i := 0; i < len(result); i++ {
					if result[i].Error != nil {
						resultList = append(resultList, result[i].Error)
					} else {
						resultList = append(resultList, result[i].Result)
					}
				}
				return output.Render(os.Stdout, config, output.Marshaled{Value: resultList})
			}
			if result[0].Error != nil {
				return output.Render(os.Stdout, config, output.Marshaled{Value: result[0].Error})
			}
			return output.Render(os.Stdout, config, output.Marshaled{Value: result[0].Result})
		}
	} else if cCmd := child.MakeExecuteCommand(); cCmd != nil && child.Cobra().RunE == nil {
		child.Cobra().RunE = func(_ *cobra.Command, args []string) error {
			if loader := child.ConfigLoader(); loader != nil {
				if err := loader(config); err != nil {
					return fmt.Errorf("Config load: %v", err)
				}
			}
			response, err := cCmd(args)
			if err != nil {
				return err
			}
			if !config.OutputHuman() {
				return handleOutput(response, config.Output())
			}
			return child.HandleOutput(os.Stdout, response)
		}
	}

	if cCmd := child.MakePreExecuteCommand(); cCmd != nil && child.Cobra().PreRunE == nil {
		child.Cobra().PreRunE = func(_ *cobra.Command, args []string) error {
			if loader := child.ConfigLoader(); loader != nil {
				if err := loader(config); err != nil {
					return fmt.Errorf("Config load: %v", err)
				}
			}
			return cCmd(args)
		}
	}
	if cCmd := child.MakePersistentPreExecuteCommand(); cCmd != nil && child.Cobra().PersistentPreRunE == nil {
		child.Cobra().PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if loader := child.ConfigLoader(); loader != nil {
				if err := loader(config); err != nil {
					return fmt.Errorf("Config load: %v", err)
				}
			}
			return cCmd(args)
		}
	}

	curHelp := child.Cobra().HelpFunc()
	child.Cobra().SetHelpFunc(func(cCmd *cobra.Command, args []string) {
		for cmd := child; cmd != nil; cmd = cmd.Parent() {
			for _, v := range cmd.Config().BoundFlags() {
				if !cmd.Config().IsSet(v.Name) {
					continue
				}
				v.DefValue = cmd.Config().GetString(v.Name)
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

// BaseCommand is the base type for all commands, implementing Command
type BaseCommand struct {
	cobra            *cobra.Command
	name             string
	parent           Command
	childrenPos      map[Command]int
	nextChildSortPos int
	config           *config.Config
	configLoader     func(config *config.Config) error
}

// Name returns the name of the command
func (s *BaseCommand) Name() string {
	return s.name
}

// SetConfig sets the configuration used
func (s *BaseCommand) SetConfig(config *config.Config) {
	s.config = config
}

// SetChild sets command as the child of this command
func (s *BaseCommand) SetChild(command Command) {
	if command == nil {
		return
	}
	if _, alreadyChild := s.childrenPos[command]; alreadyChild {
		return
	}
	if s.childrenPos == nil {
		s.childrenPos = make(map[Command]int)
	}
	s.childrenPos[command] = s.nextChildSortPos
	s.nextChildSortPos++
	if command.Parent() != s {
		command.SetParent(s)
	}
}

// DeleteChild removes command from the children of this command
func (s *BaseCommand) DeleteChild(command Command) {
	if command.Parent() == s {
		command.SetParent(nil)
	}
	delete(s.childrenPos, command)
}

// Children returns a list of all the child commands of this command (eg. including the children of children)
func (s *BaseCommand) Children() []Command {
	var (
		r      []Command
		sorted []Command
	)
	for child := range s.childrenPos {
		sorted = append(sorted, child)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return s.childrenPos[sorted[i]] < s.childrenPos[sorted[j]]
	})
	for _, child := range sorted {
		r = append(r, child)
		r = append(r, child.Children()...)
	}
	return r
}

// SetParent sets the parent of this command to given command
func (s *BaseCommand) SetParent(command Command) {
	if s.parent != nil {
		s.DeleteChild(command)
	}
	s.parent = command
	if s.parent != nil {
		s.parent.SetChild(s)
	}
}

// Parent returns the parent of the command
func (s *BaseCommand) Parent() Command {
	return s.parent
}

// SetFlags parses the given flags
func (s *BaseCommand) SetFlags(flags []string) error {
	return s.Cobra().Flags().Parse(flags)
}

// InitCommand can be overriden to handle flag registration.
// A hook to handle flag registration.
// The config values are not available during this hook. Register a cobra hook to use them. You can set defaults though.
func (s *BaseCommand) InitCommand() {
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

// Namespace returns the namespace of this command from the chain of parent commands
// The format is cmdRoot.child1.child2.childN
// No namespace is returned for the root command (parent == nil)
func (s *BaseCommand) Namespace() string {
	var (
		sb    strings.Builder
		names []string
	)
	for c := s.parent; c != nil; c = c.Parent() {
		// Skip root command name in namespace
		if c.Parent() == nil {
			continue
		}
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

// Cobra returns the underlying *cobra.Command
func (s *BaseCommand) Cobra() *cobra.Command {
	return s.cobra
}

// Config //

// Config implements Command.Config, returning the *config.Config of the command
func (s *BaseCommand) Config() *config.Config {
	return s.config
}

// SetConfigLoader implements Command.SetConfigLoader, setting internal config loader
func (s *BaseCommand) SetConfigLoader(fn func(config *config.Config) error) {
	s.configLoader = fn
}

// ConfigLoader implements Command.ConfigLoader, returning the specified ConfigLoader
func (s *BaseCommand) ConfigLoader() func(config *config.Config) error {
	return s.configLoader
}

// Flags //

// AddFlags adds a flagset to the command and binds config value into it with namespace
func (s *BaseCommand) AddFlags(flags *pflag.FlagSet) {
	if flags == nil {
		panic("Nil flagset")
	}
	flags.VisitAll(func(flag *pflag.Flag) {
		s.Cobra().Flags().AddFlag(flag)
	})
	s.config.ConfigBindFlagSet(flags)
}

// AddPersistentFlags adds a persistent flagset to the command and binds config value into it with namespace
func (s *BaseCommand) AddPersistentFlags(flags *pflag.FlagSet) {
	if flags == nil {
		panic("Nil flagset")
	}
	flags.VisitAll(func(flag *pflag.Flag) {
		s.Cobra().PersistentFlags().AddFlag(flag)
	})
	s.config.ConfigBindFlagSet(flags)
}

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
		s.cobra.Use = s.name
		return
	}
	s.cobra.Use = fmt.Sprintf("%s %s", s.name, help)
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

// HandleOutput should be overwritten by the actual command implementations
// It is expected to write output, given in out, to writer
func (s *BaseCommand) HandleOutput(io.Writer, interface{}) error {
	return nil
}

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
