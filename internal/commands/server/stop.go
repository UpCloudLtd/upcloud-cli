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

// StopCommand creates the "server stop" command
func StopCommand() commands.Command {
	return &stopCommand{
		BaseCommand: commands.New("stop", "Stop a server"),
	}
}

type stopCommand struct {
	*commands.BaseCommand
	params stopParams
}

type stopParams struct {
	request.StopServerRequest
	timeout int
}

var defaultStopParams = &stopParams{
	StopServerRequest: request.StopServerRequest{
		StopType: upcloud.StopTypeSoft,
	},
	timeout: 120,
}

// InitCommand implements Command.InitCommand
func (s *stopCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.Config().Service.(service.Server)))

	flags := &pflag.FlagSet{}
	flags.StringVar(&s.params.StopType, "type", defaultStopParams.StopType, "The type of stop operation. Available: soft, hard")
	flags.IntVar(&s.params.timeout, "timeout", defaultStartParams.timeout, "Stop timeout in seconds. Available: 1-600")
	s.AddFlags(flags)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *stopCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		timeout, err := time.ParseDuration(strconv.Itoa(s.params.timeout) + "s")
		if err != nil {
			return nil, err
		}

		s.params.Timeout = timeout
		svc := s.Config().Service.(service.Server)

		return Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.StopServerRequest
				req.UUID = uuid
				return &req
			},
			Service: svc,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.StopServerRequest).UUID },
				InteractiveUI: s.Config().InteractiveUI(),
				WaitMsg:       "shutdown request sent",
				WaitFn:        waitForServer(svc, upcloud.ServerStateStopped, s.Config().ClientTimeout()),
				MaxActions:    maxServerActions,
				ActionMsg:     "Stopping",
				Action: func(req interface{}) (interface{}, error) {
					return svc.StopServer(req.(*request.StopServerRequest))
				},
			},
		}.Send(args)
	}
}
