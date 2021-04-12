package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// New returns a BaseCommand that implements Command. It is used as a base to create custom commands from.
func New(name, usage string) *BaseCommand {
	return &BaseCommand{
		name:  name,
		cobra: &cobra.Command{Use: name, Short: usage},
	}
}

// Command is the base command type for all commands.
type Command interface {
	InitCommand()
	CobraCommand
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

func commandNoneRunE(nc NoArgumentCommand, config *config.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		svc, err := config.CreateService()
		if err != nil {
			return fmt.Errorf("cannot create service: %w", err)
		}

		executor := NewExecutor(config, svc)
		res, err := nc.ExecuteWithoutArguments(executor)
		if err != nil {
			return output.Render(os.Stdout, config, output.Marshaled{Value: err})
		}
		return output.Render(os.Stdout, config, res)
	}
}

func commandSingleRunE(nc SingleArgumentCommand, config *config.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 || args[0] == "" {
			return fmt.Errorf("exactly 1 argument is required")
		}
		svc, err := config.CreateService()
		if err != nil {
			return fmt.Errorf("cannot create service: %w", err)
		}

		executor := NewExecutor(config, svc)

		var argumentResolver resolver.Resolver
		if resolve, ok := nc.(resolver.ResolutionProvider); ok {
			argumentResolver, err = resolve.Get(svc)
			if err != nil {
				return fmt.Errorf("cannot create resolver: %w", err)
			}
		}
		var executeArgument = args[0]
		if argumentResolver != nil {
			if res, err := argumentResolver(args[0]); err == nil {
				executeArgument = res
			} else {
				executor.NewLogEntry(fmt.Sprintf("invalid argument: %v", err))
				return fmt.Errorf("cannot map argument '%v': %w", args[0], err)
			}
		}
		res, err := nc.ExecuteSingleArgument(executor, executeArgument)
		if err != nil {
			return output.Render(os.Stdout, config, output.Marshaled{Value: err})
		}
		return output.Render(os.Stdout, config, res)
	}
}

func commandMultiRunE(nc MultipleArgumentCommand, config *config.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || args[0] == "" {
			return fmt.Errorf("at least one argument is required")
		}
		svc, err := config.CreateService()
		if err != nil {
			return fmt.Errorf("cannot create service: %w", err)
		}

		executor := NewExecutor(config, svc)

		var argumentResolver resolver.Resolver
		if resolve, ok := nc.(resolver.ResolutionProvider); ok {
			argumentResolver, err = resolve.Get(svc)
			if err != nil {
				return fmt.Errorf("cannot create resolver: %w", err)
			}
		}

		returnChan := make(chan executeResult)
		workerCount := nc.MaximumExecutions()
		workerQueue := make(chan int, workerCount)

		// push initial workers into the worker queue
		for n := 0; n < workerCount; n++ {
			workerQueue <- n
		}

		// make a copy of the original args to pass into the workers
		argQueue := args
		if len(argQueue) == 0 {
			// no argument commands *still* need to run so trigger a single execution with "" as the argument
			argQueue = []string{""}
		}

		results := make([]executeResult, 0, len(args))
		renderTicker := time.NewTicker(100 * time.Millisecond)

		for {
			select {
			case workerID := <-workerQueue:
				// got an idle worker
				if len(argQueue) == 0 {
					// we are out of arguments to process, just let the worker exit
					break
				}
				arg := argQueue[0]
				argQueue = argQueue[1:]
				// trigger execution in a goroutine
				go func(index int, argument string) {
					defer func() {
						// return worker to queue when exiting
						workerQueue <- workerID
					}()
					executeArgument := argument
					if argumentResolver != nil {
						if res, err := argumentResolver(argument); err == nil {
							executeArgument = res
						} else {
							executor.NewLogEntry(fmt.Sprintf("invalid argument: %v", err))
							returnChan <- executeResult{Job: index, Error: fmt.Errorf("cannot map argument '%v': %w", argument, err)}
							return
						}
					}
					res, err := nc.Execute(executor, executeArgument)
					// return result
					returnChan <- executeResult{Job: index, Result: res, Error: err}
				}(workerID, arg)
			case res := <-returnChan:
				// got a result from a worker
				results = append(results, res)
				if len(results) >= len(args) {
					// we're done, update ui for the last time and render the results
					executor.Update()
					if len(results) > 1 {
						resultList := make([]output.Output, len(results))
						for i := 0; i < len(results); i++ {
							if results[i].Error != nil {
								resultList[i] = output.Error{Value: results[i].Error}
							} else {
								resultList[i] = results[i].Result
							}
						}
						return output.Render(os.Stdout, config, resultList...)
					}

					if results[0].Error != nil {
						return output.Render(os.Stdout, config, output.Error{Value: results[0].Error})
					}
					return output.Render(os.Stdout, config, results[0].Result)
				}
			case <-renderTicker.C:
				if config.InteractiveUI() {
					executor.Update()
				}
			}
		}
	}
}

// BuildCommand sets up a Command with the specified config and adds it to Cobra
func BuildCommand(child Command, parent *cobra.Command, config *config.Config) Command {
	child.Cobra().Flags().SortFlags = false

	// Need to set child command in the end as otherwise HelpFunc() returns the parent's helpfunc
	if parent != nil {
		parent.AddCommand(child.Cobra())
	}

	// TODO: taken out, do we need this?
	// if nsCmd, ok := child.(namespace); ok {
	//   config.SetNamespace(child.Namespace())
	// }

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
	switch typedChild := child.(type) {
	case NoArgumentCommand:
		child.Cobra().RunE = commandNoneRunE(typedChild, config)
	case SingleArgumentCommand:
		child.Cobra().RunE = commandSingleRunE(typedChild, config)
	case MultipleArgumentCommand:
		child.Cobra().RunE = commandMultiRunE(typedChild, config)
	default:
		// no execution found on this command, eg. most likely an 'organizational' command
		// so just show usage
		child.Cobra().RunE = func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		}
	}

	// Apply viper value to the help
	curHelp := child.Cobra().HelpFunc()
	child.Cobra().SetHelpFunc(func(cCmd *cobra.Command, args []string) {
		child.Cobra().Flags().VisitAll(func(f *pflag.Flag) {
			// TODO: reimplement
			/*config.SetNamespace(child.Namespace())

			if !child.Config().IsSet(f.Name) {
				return
			}
			f.DefValue = child.Config().GetString(f.Name)*/
		})
		curHelp(cCmd, args)
	})

	return child
}

// BaseCommand is the base type for all commands, implementing Command
type BaseCommand struct {
	cobra *cobra.Command
	name  string
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
	//	s.config.ConfigBindFlagSet(flags)
}

// InitCommand can be overriden to handle flag registration.
// A hook to handle flag registration.
// The config values are not available during this hook. Register a cobra hook to use them. You can set defaults though.
func (s *BaseCommand) InitCommand() {
}

// Cobra returns the underlying *cobra.Command
func (s *BaseCommand) Cobra() *cobra.Command {
	return s.cobra
}
