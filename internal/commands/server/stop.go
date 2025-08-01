package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// StopCommand creates the "server stop" command
func StopCommand() commands.Command {
	return &stopCommand{
		BaseCommand: commands.New(
			"stop",
			"Stop a server",
			"upctl server stop 00cbe2f3-4cf9-408b-afee-bd340e13cdd8",
			"upctl server stop 00cbe2f3-4cf9-408b-afee-bd340e13cdd8 0053a6f5-e6d1-4b0b-b9dc-b90d0894e8d0",
			"upctl server stop my_server",
			"upctl server stop --wait my_server",
		),
	}
}

type stopCommand struct {
	*commands.BaseCommand
	StopType string
	wait     config.OptionalBoolean
	resolver.CachingServer
	completion.StartedServer
}

// InitCommand implements Command.InitCommand
func (s *stopCommand) InitCommand() {
	// XXX: findout what to do with risky params (timeout actions)
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.StopType, "type", defaultStopType, "The type of stop operation. Available: soft, hard")
	config.AddToggleFlag(flags, &s.wait, "wait", false, "Wait for server to be in stopped state before returning.")
	s.AddFlags(flags)
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("type", cobra.FixedCompletions(stopTypes, cobra.ShellCompDirectiveNoFileComp)))
}

func stop(exec commands.Executor, uuid, stopType string, wait bool) (output.Output, error) {
	svc := exec.Server()
	msg := fmt.Sprintf("Stopping server %v", uuid)
	exec.PushProgressStarted(msg)

	res, err := svc.StopServer(exec.Context(), &request.StopServerRequest{
		UUID:     uuid,
		StopType: stopType,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if wait {
		waitForServerState(uuid, upcloud.ServerStateStopped, exec, msg)
	} else {
		exec.PushProgressSuccess(msg)
	}

	return output.OnlyMarshaled{Value: res}, nil
}

// Execute implements commands.MultipleArgumentCommand
func (s *stopCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	return stop(exec, uuid, s.StopType, s.wait.Value())
}
