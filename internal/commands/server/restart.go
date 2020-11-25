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

var DefaultRestartParams = &restartParams{
	RestartServerRequest: request.RestartServerRequest{
		StopType:      "soft",
		TimeoutAction: "ignore",
	},
	timeout: 120,
}

func (s *restartCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.service))

	s.params = restartParams{RestartServerRequest: request.RestartServerRequest{}}
	flags := &pflag.FlagSet{}

	flags.StringVar(&s.params.StopType, "stop-type", DefaultRestartParams.StopType, "Restart type\nAvailable: soft, hard")
	flags.StringVar(&s.params.TimeoutAction, "timeout-action", DefaultRestartParams.TimeoutAction, "Action to take if timeout limit is exceeded\nAvailable: destroy, ignore")
	flags.IntVar(&s.params.timeout, "timeout", DefaultRestartParams.timeout, "Stop timeout in seconds\nAvailable: 1-600")
	flags.IntVar(&s.params.Host, "host", DefaultRestartParams.Host, "Use this to restart the VM on a specific host. Refers to value from host attribute. Only available for private cloud hosts")

	s.AddFlags(flags)
}

func (s *restartCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		timeout, err := time.ParseDuration(strconv.Itoa(s.params.timeout) + "s")
		if err != nil {
			return nil, err
		}
		s.params.Timeout = timeout

		return Request{
			BuildRequest: func(server *upcloud.Server) interface{} {
				req := s.params.RestartServerRequest
				req.UUID = server.UUID
				return &req
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.RestartServerRequest).UUID },
				InteractiveUI: s.Config().InteractiveUI(),
				WaitMsg:       "restart request sent",
				WaitFn:        WaitForServerFn(s.service, upcloud.ServerStateStarted, s.Config().ClientTimeout()),
				MaxActions:    maxServerActions,
				ActionMsg:     "Restarting",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.RestartServer(req.(*request.RestartServerRequest))
				},
			},
		}.Send(args)
	}
}
