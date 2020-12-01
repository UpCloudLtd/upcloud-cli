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

func StartCommand(service service.Server) commands.Command {
	return &startCommand{
		BaseCommand: commands.New("start", "Start a server"),
		service:     service,
	}
}

type startCommand struct {
	*commands.BaseCommand
	service service.Server
	params  startParams
}

type startParams struct {
	request.StartServerRequest
	timeout int
}

var DefaultStartParams = &startParams{
	StartServerRequest: request.StartServerRequest{},
	timeout:            120,
}

func (s *startCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.service))

	flags := &pflag.FlagSet{}
	flags.IntVar(&s.params.AvoidHost, "avoid-host", DefaultStartParams.AvoidHost, "Avoid specific host when starting a server")
	flags.IntVar(&s.params.Host, "host", DefaultStartParams.Host, "Start server on a specific host. Note that this is generally available for private clouds only")
	flags.IntVar(&s.params.timeout, "timeout", DefaultStartParams.timeout, "Stop timeout in seconds\nAvailable: 1-600")
	s.AddFlags(flags)
}

func (s *startCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		timeout, err := time.ParseDuration(strconv.Itoa(s.params.timeout) + "s")
		if err != nil {
			return nil, err
		}
		s.params.Timeout = timeout

		return Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.StartServerRequest
				req.UUID = uuid
				return &req
			},
			Service: s.service,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.StartServerRequest).UUID },
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    maxServerActions,
				WaitMsg:       "starting server",
				WaitFn:        WaitForServerFn(s.service, upcloud.ServerStateStarted, s.Config().ClientTimeout()),
				ActionMsg:     "Starting",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.StartServer(req.(*request.StartServerRequest))
				},
			},
		}.Send(args)
	}
}
