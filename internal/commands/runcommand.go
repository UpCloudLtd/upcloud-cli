package commands

import (
	"fmt"
	"io"

	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	internal "github.com/UpCloudLtd/upcloud-cli/internal/service"

	"github.com/gemalto/flume"
)

var logger = flume.New("runcommand")

func commandRunE(command Command, service internal.AllServices, config *config.Config, args []string) error {
	// Cobra validations were successful
	command.Cobra().SilenceUsage = true

	cmdLogger := logger.With("command", command.Cobra().CommandPath())
	executor := NewExecutor(config, service, cmdLogger)

	var w io.Writer

	if config.OutputHuman() {
		w = command.Cobra().OutOrStderr()
	} else {
		w = command.Cobra().OutOrStdout()
	}

	switch typedCommand := command.(type) {
	case NoArgumentCommand:
		cmdLogger.Debug("executing without arguments", "arguments", args)
		// need to pass in fake arguments here, to actually trigger execution
		results, err := execute(typedCommand, executor, []string{""}, 1,
			// FIXME: this bit panics go-critic unlambda check, figure out why and report upstream
			func(exec Executor, fake string) (output.Output, error) {
				return typedCommand.ExecuteWithoutArguments(exec)
			})
		if err != nil {
			return err
		}
		return output.Render(w, config.Output(), results...)
	case SingleArgumentCommand:
		cmdLogger.Debug("executing single argument", "arguments", args)
		// make sure we have an argument
		if len(args) != 1 || args[0] == "" {
			return fmt.Errorf("exactly one positional argument is required")
		}
		results, err := execute(typedCommand, executor, args, 1, typedCommand.ExecuteSingleArgument)
		if err != nil {
			return err
		}
		return output.Render(w, config.Output(), results...)
	case MultipleArgumentCommand:
		cmdLogger.Debug("executing multi argument", "arguments", args)
		// make sure we have arguments
		if len(args) < 1 {
			return fmt.Errorf("at least one positional argument is required")
		}
		results, err := execute(typedCommand, executor, args, typedCommand.MaximumExecutions(), typedCommand.Execute)
		if err != nil {
			return err
		}
		return output.Render(w, config.Output(), results...)
	default:
		// no execution found on this command, eg. most likely an 'organizational' command
		// so just show usage
		cmdLogger.Debug("no execution found", "arguments", args)
		return command.Cobra().Usage()
	}
}

type resolvedArgument struct {
	Resolved string
	Error    error
	Original string
}

func resolveArguments(nc Command, svc internal.AllServices, args []string) (out []resolvedArgument, err error) {
	if resolve, ok := nc.(resolver.ResolutionProvider); ok {
		argumentResolver, err := resolve.Get(svc)
		if err != nil {
			return nil, fmt.Errorf("cannot get resolver: %w", err)
		}
		for _, arg := range args {
			resolved, err := argumentResolver(arg)
			out = append(out, resolvedArgument{Resolved: resolved, Error: err, Original: arg})
		}
	} else {
		for _, arg := range args {
			out = append(out, resolvedArgument{Resolved: arg, Original: arg})
		}
	}
	return
}

func execute(command Command, executor Executor, args []string, parallelRuns int, executeCommand func(exec Executor, arg string) (output.Output, error)) ([]output.Output, error) {
	resolvedArgs, err := resolveArguments(command, executor.All(), args)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve command line arguments: %w", err)
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

	outputs := make([]output.Output, 0, len(args))
	executor.Debug("starting work", "workers", workerCount)
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
					executor.Debug("worker exiting", "worker", index)
					workerQueue <- workerID
				}()
				if argument.Error != nil {
					// argument wasn't parsed correctly, pass the error on
					executor.Debug("worker got invalid argument", "worker", index, "error", argument.Error)
					returnChan <- executeResult{Job: index, Error: fmt.Errorf("cannot resolve argument: %w", argument.Error)}
				} else {
					// otherwise, execute and return results
					executor.Debug("worker starting", "worker", index, "argument", argument.Resolved)
					res, err := executeCommand(executor.WithLogger("worker", index, "argument", argument.Resolved), argument.Resolved)
					returnChan <- executeResult{Job: index, Result: res, Error: err, ResolvedArgument: argument}
				}
			}(workerID, arg)
		case res := <-returnChan:
			// got a result from a worker
			if res.Error != nil {
				outputs = append(outputs, output.Error{
					Value:    res.Error,
					Resolved: res.ResolvedArgument.Resolved,
					Original: res.ResolvedArgument.Original,
				})
			} else {
				outputs = append(outputs, res.Result)
			}

			if len(outputs) >= len(args) {
				executor.Debug("execute done")
				// We're done, update ui for the last time and render the results
				executor.StopProgressLog()
				return outputs, nil
			}
		}
	}
}
