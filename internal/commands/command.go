package commands

import (
	"encoding/json"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/mapper"
	"github.com/UpCloudLtd/cli/internal/output"
	"io"
	"os"
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
	SetFlags(flags []string) error
	Name() string
	InitCommand()
	SetupFlags()
	MakeExecuteCommand() func(args []string) (interface{}, error)
	Config() *config.Config
	SetConfig(*config.Config)
	HandleOutput(writer io.Writer, out interface{}) error
	HandleError(err error)
	Namespace() string
	CobraCommand
}

// NewCommand is a new container for commands, currently still including the old interface until we can deprecate it
type NewCommand interface {
	Execute(exec Executor, arg string) (output.Command, error)
	MaximumExecutions() int
	ArgumentMapper() (mapper.Argument, error)
	Command
}

// CobraCommand is an interface for commands that can refer back to their base cobra.Command
type CobraCommand interface {
	Cobra() *cobra.Command
}

// type namespace interface {
// }

// BuildCommand sets up a Command with the specified config and adds it to Cobra
func BuildCommand(child Command, parent *cobra.Command, config *config.Config) Command {
	child.SetConfig(config)
	child.Cobra().Flags().SortFlags = false

	// Need to set child command in the end as otherwise HelpFunc() returns the parent's helpfunc
	if parent != nil {
		parent.AddCommand(child.Cobra())
	}

	// if nsCmd, ok := child.(namespace); ok {
	config.SetNamespace(child.Namespace())
	// }

	// Init
	child.InitCommand()

	// Run
	if nc, ok := child.(NewCommand); ok {
		child.Cobra().RunE = func(cmd *cobra.Command, args []string) error {
			svc, err := config.CreateService()
			if err != nil {
				return fmt.Errorf("cannot create service: %w", err)
			}

			executor := NewExecutor(config, svc)
			argmapper, err := nc.ArgumentMapper()
			if err != nil {
				return fmt.Errorf("invalid mapper: %w", err)
			}

			returnChan := make(chan executeResult)
			if len(args) > 0 {
				for i, arg := range args {
					go func(index int, argument string) {
						executeArgument := argument
						if argmapper != nil {
							if res, err := argmapper(argument); err == nil {
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
			} else { // no args cmds
				go func() {
					res, err := nc.Execute(executor, "")
					returnChan <- executeResult{Job: 0, Result: res, Error: err}
				}()
			}

			result := map[int]executeResult{}
			renderTicker := time.NewTicker(100 * time.Millisecond)
		waitForResults:
			for {
				select {
				case res := <-returnChan:
					result[res.Job] = res
					if len(result) >= len(args) {
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
			// calling this just to init service in non-refactored commands as well, TODO: remove when refactored..
			if _, err := config.CreateService(); err != nil {
				return fmt.Errorf("cannot setup service client: %w", err)
			}

			// Apply values set from viper if any
			child.Cobra().Flags().VisitAll(func(f *pflag.Flag) {
				config.SetNamespace(child.Namespace())
				if !child.Config().IsSet(f.Name) {
					return
				}
				if f.Changed {
					return
				}
				if err := f.Value.Set(child.Config().GetString(f.Name)); err != nil {
					return
				}
			})

			// Run
			response, err := cCmd(args)
			if err != nil {
				return err
			}

			// Handle output
			if !config.OutputHuman() {
				return handleOutput(response, config.Output())
			}
			return child.HandleOutput(os.Stdout, response)
		}
	}

	// Apply viper value to the help
	curHelp := child.Cobra().HelpFunc()
	child.Cobra().SetHelpFunc(func(cCmd *cobra.Command, args []string) {
		child.Cobra().Flags().VisitAll(func(f *pflag.Flag) {
			config.SetNamespace(child.Namespace())

			if !child.Config().IsSet(f.Name) {
				return
			}
			f.DefValue = child.Config().GetString(f.Name)
		})
		curHelp(cCmd, args)
	})

	return child
}

// BaseCommand is the base type for all commands, implementing Command
type BaseCommand struct {
	cobra  *cobra.Command
	name   string
	config *config.Config
}

// Name returns the name of the command
func (s *BaseCommand) Name() string {
	return s.name
}

// SetFlags parses the given flags
func (s *BaseCommand) SetFlags(flags []string) error {
	return s.Cobra().Flags().Parse(flags)
}

// SetupFlags allows adding defined flags to cobra cmd
func (s *BaseCommand) SetupFlags() {
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

// Namespace returns the namespace of this command from the chain of parent commands
// The format is cmdRoot.child1.child2.childN
// No namespace is returned for the root command (parent == nil)
func (s *BaseCommand) Namespace() string {
	var (
		sb    strings.Builder
		names []string
	)
	for c := s.cobra; c != nil; c = c.Parent() {
		// Skip root command name in namespace
		if c.Parent() == nil {
			continue
		}
		names = append(names, c.Name())
	}

	for i := len(names) - 1; i >= 0; i-- {
		sb.WriteString(names[i])
		if i != 0 {
			sb.WriteString(".")
		}
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

// SetConfig sets the configuration used
func (s *BaseCommand) SetConfig(config *config.Config) {
	s.config = config
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
