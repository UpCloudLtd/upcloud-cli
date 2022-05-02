package commands

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// New returns a BaseCommand that implements Command. It is used as a base to create custom commands from.
func New(name, usage string, examples ...string) *BaseCommand {
	return &BaseCommand{
		cobra: &cobra.Command{
			Use:     name,
			Short:   usage,
			Example: strings.Join(examples, "\n"),
		},
	}
}

// Command is the base command type for all commands.
type Command interface {
	InitCommand()
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

	// Set up completion, if necessary
	if cp, ok := child.(completion.Provider); ok {
		child.Cobra().ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			svc, err := config.CreateService()
			if err != nil {
				// TODO: debug log fmt.Sprintf("cannot create service for completion: %v", err)
				return completion.None(toComplete)
			}
			return cp.CompleteArgument(svc, toComplete)
		}
	}
	if rp, ok := child.(resolver.ResolutionProvider); ok {
		child.Cobra().Use = fmt.Sprintf("%s %s", child.Cobra().Name(), rp.PositionalArgumentHelp())
	} else {
		// not sure if we really need this?
		child.Cobra().Use = child.Cobra().Name()
	}

	// Set run
	child.Cobra().RunE = func(cmd *cobra.Command, args []string) error {
		// Do not create service for offline commands, e.g. upctl version
		if _, ok := child.(OfflineCommand); ok {
			return commandRunE(child, nil, config, args)
		}

		service, err := config.CreateService()
		if err != nil {
			// Error was caused by missing credentials, not incorrect command
			child.Cobra().SilenceUsage = true

			return fmt.Errorf("cannot create service: %w", err)
		}
		return commandRunE(child, service, config, args)
	}

	return child
}

// BaseCommand is the base type for all commands, implementing Command
type BaseCommand struct {
	cobra *cobra.Command
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
	})
}

// InitCommand can be overridden to handle flag registration.
// A hook to handle flag registration.
// The config values are not available during this hook. Register a cobra hook to use them. You can set defaults though.
func (s *BaseCommand) InitCommand() {
}

// Cobra returns the underlying *cobra.Command
func (s *BaseCommand) Cobra() *cobra.Command {
	return s.cobra
}
