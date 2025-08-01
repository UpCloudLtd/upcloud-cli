package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// StartCommand creates the "server start" command
func StartCommand() commands.Command {
	return &startCommand{
		BaseCommand: commands.New(
			"start",
			"Start a server",
			"upctl server start 00038afc-d526-4148-af0e-d2f1eeaded9b",
			"upctl server start 00038afc-d526-4148-af0e-d2f1eeaded9b 0053a6f5-e6d1-4b0b-b9dc-b90d0894e8d0",
			"upctl server start my_server1",
		),
	}
}

type startCommand struct {
	*commands.BaseCommand
	completion.StoppedServer
	resolver.CachingServer
	host      int
	avoidHost int
}

// InitCommand implements Command.InitCommand
func (s *startCommand) InitCommand() {
	fs := &pflag.FlagSet{}

	fs.IntVar(&s.avoidHost, "avoid-host", 0, avoidHostDescription)
	fs.IntVar(&s.host, "host", 0, hostDescription)

	s.AddFlags(fs)
}

func (s *startCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("avoid-host", namedargs.CompletionFunc(completion.HostID{}, cfg)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("host", namedargs.CompletionFunc(completion.HostID{}, cfg)))
}

// Execute implements commands.MultipleArgumentCommand
func (s *startCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Server()
	msg := fmt.Sprintf("Starting server %v", uuid)
	exec.PushProgressStarted(msg)

	res, err := svc.StartServer(exec.Context(), &request.StartServerRequest{
		UUID:      uuid,
		AvoidHost: s.avoidHost,
		Host:      s.host,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
