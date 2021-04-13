package commands

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	internal "github.com/UpCloudLtd/cli/internal/service"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func commandRunE(command Command, config *config.Config) func(cmd *cobra.Command, args []string) error {

	switch typedCommand := command.(type) {
	case NoArgumentCommand:
		return func(cmd *cobra.Command, args []string) error {
			// need to pass in fake arguments here, to actually trigger execution
			return execute(typedCommand, config, []string{""}, 1,
				// FIXME: these bits panic go-critic unlambda check, figure out why and report upstream
				func(exec Executor, fake string) (output.Output, error) {
					return typedCommand.ExecuteWithoutArguments(exec)
				})
		}
	case SingleArgumentCommand:
		return func(cmd *cobra.Command, args []string) error {
			// make sure we have an argument
			if len(args) != 1 || args[0] == "" {
				return fmt.Errorf("exactly 1 argument is required")
			}
			// FIXME: these bits panic go-critic unlambda check, figure out why and report upstream
			return execute(typedCommand, config, args, 1, func(exec Executor, arg string) (output.Output, error) {
				return typedCommand.ExecuteSingleArgument(exec, arg)
			})
		}
	case MultipleArgumentCommand:
		return func(cmd *cobra.Command, args []string) error {
			// make sure we have arguments
			if len(args) < 1 {
				return fmt.Errorf("at least one argument is required")
			}
			// FIXME: these bits panic go-critic unlambda check, figure out why and report upstream
			return execute(typedCommand, config, args, typedCommand.MaximumExecutions(), func(exec Executor, arg string) (output.Output, error) {
				return typedCommand.Execute(exec, arg)
			})
		}
	default:
		// no execution found on this command, eg. most likely an 'organizational' command
		// so just show usage
		return func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		}
	}

}

func render(config *config.Config, results []executeResult) error {
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

type resolvedArgument struct {
	Resolved string
	Error    error
}

func resolveArguments(nc Command, svc internal.AllServices, args []string) (out []resolvedArgument, err error) {
	if resolve, ok := nc.(resolver.ResolutionProvider); ok {
		argumentResolver, err := resolve.Get(svc)
		if err != nil {
			return nil, fmt.Errorf("cannot create resolver: %w", err)
		}
		for _, arg := range args {
			resolved, err := argumentResolver(arg)
			out = append(out, resolvedArgument{Resolved: resolved, Error: err})
		}
	} else {
		for _, arg := range args {
			out = append(out, resolvedArgument{Resolved: arg})
		}
	}
	return
}

func execute(command Command, config *config.Config, args []string, parallelRuns int, executeCommand func(exec Executor, arg string) (output.Output, error)) error {
	svc, err := config.CreateService()
	if err != nil {
		return fmt.Errorf("cannot create service: %w", err)
	}
	executor := NewExecutor(config, svc)
	resolvedArgs, err := resolveArguments(command, svc, args)
	if err != nil {
		return fmt.Errorf("cannot create resolver: %w", err)
	}

	returnChan := make(chan executeResult)
	workerCount := parallelRuns
	workerQueue := make(chan int, workerCount)

	// push initial workers into the worker queue
	for n := 0; n < workerCount; n++ {
		workerQueue <- n
	}

	// make a copy of the original args to pass into the workers
	argQueue := resolvedArgs

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
			go func(index int, argument resolvedArgument) {
				defer func() {
					// return worker to queue when exiting
					workerQueue <- workerID
				}()
				if argument.Error != nil {
					// argument wasn't parsed correctly, pass the error on
					returnChan <- executeResult{Job: index, Error: fmt.Errorf("cannot resolve argument: %w", argument.Error)}
				} else {
					// otherwise, execute and return results
					res, err := executeCommand(executor, argument.Resolved)
					returnChan <- executeResult{Job: index, Result: res, Error: err}
				}
			}(workerID, arg)
		case res := <-returnChan:
			// got a result from a worker
			results = append(results, res)
			if len(results) >= len(args) {
				// we're done, update ui for the last time and render the results
				executor.Update()
				return render(config, results)
			}
		case <-renderTicker.C:
			executor.Update()
		}
	}
}
