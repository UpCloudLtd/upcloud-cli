package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
	"strconv"
	"time"
)

// RestartCommand creates the "server restart" command
func RestartCommand(service service.Server) commands.Command {
	return &restartCommand{
		BaseCommand: commands.New("restart", "Restart a server"),
		service:     service,
	}
}

type restartCommand struct {
	*commands.BaseCommand
	service service.Server
	params  restartParams
}

type restartParams struct {
	request.RestartServerRequest
	timeout int
}

var defaultRestartParams = &restartParams{
	RestartServerRequest: request.RestartServerRequest{
		StopType:      "soft",
		TimeoutAction: "ignore",
	},
	timeout: 120,
}

// InitCommand implements Command.InitCommand
func (s *restartCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.service))

	s.params = restartParams{RestartServerRequest: request.RestartServerRequest{}}
	flags := &pflag.FlagSet{}

	flags.StringVar(&s.params.StopType, "stop-type", defaultRestartParams.StopType, "Restart type. Available: soft, hard")
	flags.StringVar(&s.params.TimeoutAction, "timeout-action", defaultRestartParams.TimeoutAction, "Action to take if timeout limit is exceeded. Available: destroy, ignore")
	flags.IntVar(&s.params.timeout, "timeout", defaultRestartParams.timeout, "Stop timeout in seconds. Available: 1-600")
	flags.IntVar(&s.params.Host, "host", defaultRestartParams.Host, "Use this to restart the VM on a specific host. Refers to value from host attribute. Only available for private cloud hosts")

	s.AddFlags(flags)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *restartCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		timeout, err := time.ParseDuration(strconv.Itoa(s.params.timeout) + "s")
		if err != nil {
			return nil, err
		}
		s.params.Timeout = timeout

		return Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.RestartServerRequest
				req.UUID = uuid
				return &req
			},
			Service: s.service,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.RestartServerRequest).UUID },
				InteractiveUI: s.Config().InteractiveUI(),
				WaitMsg:       "restart request sent",
				WaitFn:        waitForServer(s.service, upcloud.ServerStateStarted, s.Config().ClientTimeout()),
				MaxActions:    maxServerActions,
				ActionMsg:     "Restarting",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.RestartServer(req.(*request.RestartServerRequest))
				},
			},
		}.Send(args)
	}
}
