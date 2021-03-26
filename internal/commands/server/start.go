package server

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
	"strconv"
	"time"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
)

// StartCommand creates the "server start" command
func StartCommand() commands.Command {
	return &startCommand{
		BaseCommand: commands.New("start", "Start a server"),
	}
}

type startCommand struct {
	*commands.BaseCommand
	params startParams
}

type startParams struct {
	request.StartServerRequest
	timeout int
}

var defaultStartParams = &startParams{
	StartServerRequest: request.StartServerRequest{},
	timeout:            120,
}

// InitCommand implements Command.InitCommand
func (s *startCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.Config()))

	flags := &pflag.FlagSet{}
	flags.IntVar(&s.params.AvoidHost, "avoid-host", defaultStartParams.AvoidHost, "Avoid specific host when starting a server.")
	flags.IntVar(&s.params.Host, "host", defaultStartParams.Host, "Start server on a specific host. Note that this is generally available for private clouds only.")
	flags.IntVar(&s.params.timeout, "timeout", defaultStartParams.timeout, "Stop timeout in seconds. Available: 1-600")
	s.AddFlags(flags)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *startCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		timeout, err := time.ParseDuration(strconv.Itoa(s.params.timeout) + "s")
		if err != nil {
			return nil, err
		}

		s.params.Timeout = timeout
		svc := s.Config().Service.(service.Server)

		return Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.StartServerRequest
				req.UUID = uuid
				return &req
			},
			Service: svc,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.StartServerRequest).UUID },
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    maxServerActions,
				WaitMsg:       "starting server",
				WaitFn:        waitForServer(svc, upcloud.ServerStateStarted, s.Config().ClientTimeout()),
				ActionMsg:     "Starting",
				Action: func(req interface{}) (interface{}, error) {
					return svc.StartServer(req.(*request.StartServerRequest))
				},
			},
		}.Send(args)
	}
}
