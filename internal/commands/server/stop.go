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

func StopCommand(service service.Server) commands.Command {
	return &stopCommand{
		BaseCommand: commands.New("stop", "Stop a server"),
		service:     service,
	}
}

type stopCommand struct {
	*commands.BaseCommand
	service service.Server
	params  stopParams
}

type stopParams struct {
	request.StopServerRequest
	timeout int
}

var DefaultStopParams = &stopParams{
	StopServerRequest: request.StopServerRequest{
		StopType: upcloud.StopTypeSoft,
	},
	timeout: 120,
}

func (s *stopCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.service))

	flags := &pflag.FlagSet{}
	flags.StringVar(&s.params.StopType, "type", DefaultStopParams.StopType, "The type of stop operation. Soft waits for the OS to shut down cleanly while hard forcibly shuts down a server")
	s.AddFlags(flags)
	s.SetPositionalArgHelp("<uuidHostnameOrTitle ...>")
}

func (s *stopCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		timeout, err := time.ParseDuration(strconv.Itoa(s.params.timeout) + "s")
		if err != nil {
			return nil, err
		}
		s.params.Timeout = timeout

		return Request{
			BuildRequest: func(server *upcloud.Server) interface{} {
				req := s.params.StopServerRequest
				req.UUID = server.UUID
				return &req
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.StopServerRequest).UUID },
				InteractiveUI: s.Config().InteractiveUI(),
				WaitMsg:       "shutdown request sent",
				WaitFn:        WaitForServerFn(s.service, upcloud.ServerStateStopped, s.Config().ClientTimeout()),
				MaxActions:    maxServerActions,
				ActionMsg:     "Stopping",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.StopServer(req.(*request.StopServerRequest))
				},
			},
		}.Send(args)
	}
}
