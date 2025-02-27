package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CommandContextKey string

const commandKey CommandContextKey = "command"

// New returns a BaseCommand that implements Command. It is used as a base to create custom commands from.
func New(name, usage string, examples ...string) *BaseCommand {
	cmd := &cobra.Command{
		Use:     name,
		Short:   usage,
		Example: strings.Join(examples, "\n"),
	}

	// Initialize BaseCommand
	baseCmd := &BaseCommand{
		cobra: cmd,
	}

	// Store reference to itself in the context - We need this to access the command in the CobraCommand interface
	// Specifically to generate the reference documentation
	cmd.SetContext(context.WithValue(context.Background(), commandKey, baseCmd))

	return baseCmd
}

// Command is the base command type for all commands.
type Command interface {
	InitCommand()
	InitCommandWithConfig(*config.Config)
	CobraCommand
}

// OfflineCommand is a command that does not need connect to the API, e.g. upctl version.
type OfflineCommand interface {
	DoesNotUseServices()
}

// NoArgumentCommand is a command that does not care about the positional arguments.
type NoArgumentCommand interface {
	Command
	ExecuteWithoutArguments(exec Executor) (output.Output, error)
}

// SingleArgumentCommand is a command that accepts exactly one positional argument.
type SingleArgumentCommand interface {
	Command
	ExecuteSingleArgument(exec Executor, arg string) (output.Output, error)
}

// MultipleArgumentCommand is a command that can accept multiple positional arguments,
// each of which will result in a (parallel) call to Execute() with the argument.
type MultipleArgumentCommand interface {
	Command
	MaximumExecutions() int
	Execute(exec Executor, arg string) (output.Output, error)
}

// CobraCommand is an interface for commands that can refer back to their base cobra.Command
type CobraCommand interface {
	Cobra() *cobra.Command
}

// BuildCommand sets up a Command with the specified config and adds it to Cobra
func BuildCommand(child Command, parent *cobra.Command, config *config.Config) Command {
	child.Cobra().Flags().SortFlags = false

	// Need to set child command in the end as otherwise HelpFunc() returns the parent's helpfunc
	if parent != nil {
		parent.AddCommand(child.Cobra())
	}

	// XXX: Maybe put back the viper default flags value to child commands
	// params?  It is was implemented back in
	// 5ece0e1b31df5d542546d81bbf2472c2e97aadff
	// How does it work:
	// A common viper instance can be shared for all commands, each flags has
	// the format Name:
	// Parent.Child1...Childn.FlagName

	// Init
	child.InitCommand()
	child.InitCommandWithConfig(config)

	// Set up completion, if necessary
	if cp, ok := child.(completion.Provider); ok {
		child.Cobra().ValidArgsFunction = func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			svc, err := config.CreateService()
			if err != nil {
				// TODO: debug log fmt.Sprintf("cannot create service for completion: %v", err)
				return completion.None(toComplete)
			}
			return cp.CompleteArgument(config.Context(), svc, toComplete)
		}
	} else {
		// Otherwise offer no completions.
		// This, rather than file completions, is the common case for our commands.
		child.Cobra().ValidArgsFunction = cobra.NoFileCompletions
	}
	if rp, ok := child.(resolver.ResolutionProvider); ok {
		child.Cobra().Use = fmt.Sprintf("%s %s", child.Cobra().Name(), rp.PositionalArgumentHelp())
	} else {
		// not sure if we really need this?
		child.Cobra().Use = child.Cobra().Name()
	}

	// Set run
	child.Cobra().RunE = func(_ *cobra.Command, args []string) error {
		// Do not create service for offline commands, e.g. upctl version
		if _, ok := child.(OfflineCommand); ok {
			return commandRunE(child, nil, config, args)
		}

		service, err := config.CreateService()
		if err != nil {
			// Error was caused by missing credentials, not incorrect command
			child.Cobra().SilenceUsage = true
			return err
		}
		return commandRunE(child, service, config, args)
	}

	return child
}

// BaseCommand is the base type for all commands, implementing Command
type BaseCommand struct {
	cobra             *cobra.Command
	deprecatedAliases []string
}

// Aliases return non deprecated aliases
func (s *BaseCommand) Aliases() []string {
	// Get all aliases from Cobra
	allAliases := s.cobra.Aliases

	// Filter out deprecated aliases
	var filteredAliases []string
	for _, alias := range allAliases {
		if !s.isDeprecatedAlias(alias) {
			filteredAliases = append(filteredAliases, alias)
		}
	}

	return filteredAliases
}

// isDeprecatedAlias checks if an alias is deprecated
func (s *BaseCommand) isDeprecatedAlias(alias string) bool {
	for _, deprecated := range s.deprecatedAliases {
		if alias == deprecated {
			return true
		}
	}
	return false
}

func (s *BaseCommand) DeprecatedAliases() []string {
	return s.deprecatedAliases
}

func (s *BaseCommand) SetDeprecatedAliases(aliases []string) {
	s.deprecatedAliases = aliases
}

// MaximumExecutions return the max executed workers
func (s *BaseCommand) MaximumExecutions() int {
	return 1
}

// AddFlags adds a flagset to the command and binds config value into it with namespace
func (s *BaseCommand) AddFlags(flags *pflag.FlagSet) {
	if flags == nil {
		panic("Nil flagset")
	}

	flags.VisitAll(func(flag *pflag.Flag) {
		s.Cobra().Flags().AddFlag(flag)
		if _, ok := flag.Annotations[FlagAnnotationNoFileCompletions]; ok {
			Must(s.Cobra().RegisterFlagCompletionFunc(flag.Name, cobra.NoFileCompletions))
		} else if values, ok := flag.Annotations[FlagAnnotationFixedCompletions]; ok {
			Must(s.Cobra().RegisterFlagCompletionFunc(flag.Name, cobra.FixedCompletions(values, cobra.ShellCompDirectiveNoFileComp)))
		}
	})
}

// InitCommand can be overridden to handle flag registration.
// A hook to handle flag registration.
// The config values are not available during this hook. Register a cobra hook to use them. You can set defaults though.
func (s *BaseCommand) InitCommand() {
}

// InitCommandWithConfig is a hook for running additional initialisations with access to config.
// E.g., can be used to set flag completion functions.
// Note that the config might not be initialized when this function is called. Thus, config has to be used in a hook function, e.g. Cobras Command.RegisterFlagCompletionFunc.
func (s *BaseCommand) InitCommandWithConfig(*config.Config) {
}

// Cobra returns the underlying *cobra.Command
func (s *BaseCommand) Cobra() *cobra.Command {
	return s.cobra
}
