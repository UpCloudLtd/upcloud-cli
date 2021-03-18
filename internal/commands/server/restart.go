package server

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/mapper"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
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
	service       service.Server
	StopType      string
	TimeoutAction string
	Timeout       time.Duration
}

func (s *restartCommand) MaximumExecutions() int {
	return maxServerActions
}

func (s *restartCommand) ArgumentMapper() (mapper.Argument, error) {
	return mapper.CachingServer(s.service)
}

func (s *restartCommand) Execute(exec commands.Executor, arg string) (output.Command, error) {
	msg := fmt.Sprintf("restarting server %v", arg)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	res, err := s.service.RestartServer(&request.RestartServerRequest{
		UUID:          arg,
		StopType:      s.StopType,
		Timeout:       s.Timeout,
		TimeoutAction: "ignore",
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		//		logline.SetMessage(fmt.Sprintf("failed (%v)", err))
		return nil, err
	}
	// TODO: implement wait here
	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()
	return output.Marshaled{Value: res}, nil
}

func (s *restartCommand) NewParent() commands.NewCommand {
	panic("implement me")
}

const defaultStopType = "soft"
const defaultTimeout = time.Duration(120) * time.Second

// InitCommand implements Command.InitCommand
func (s *restartCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.service))
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.StopType, "stop-type", defaultStopType, "Restart type\nAvailable: soft, hard")
	// TODO: reimplement? does not seem to make sense to automagically destroy servers if restart fails..
	// flags.StringVar(&s.TimeoutAction, "timeout-action", defaultTimeoutAction, "Action to take if timeout limit is exceeded\nAvailable: destroy, ignore")
	flags.DurationVar(&s.Timeout, "timeout", defaultTimeout, "Stop timeout in Go duration string\nExamples: 100ms, 1m10s, 3h")
	// TODO: reimplement? does not seem to be in use..
	// flags.IntVar(&s.Host, "host", defaultHost, "Use this to restart the VM on a specific host. Refers to value from host attribute. Only available for private cloud hosts")
	s.AddFlags(flags)
}

/*
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
*/
