package server

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strconv"
	"sync"
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
	timeout: 0,
}

func (s *restartCommand) InitCommand() {
	s.ArgCompletion(func(toComplete string) ([]string, cobra.ShellCompDirective) {
		servers, err := s.service.GetServers()
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var vals []string
		for _, v := range servers.Servers {
			vals = append(vals, v.UUID, v.Hostname)
		}
		return commands.MatchStringPrefix(vals, toComplete, false), cobra.ShellCompDirectiveNoFileComp
	})

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

		if len(args) < 1 {
			return nil, fmt.Errorf("server title or uuid is required")
		}

		var (
			restartServers []request.RestartServerRequest
			allServers     []upcloud.Server
		)

		timeout, err := time.ParseDuration(strconv.Itoa(s.params.timeout) + "s")
		if err != nil {
			return nil, err
		}

		for _, v := range args {
			server, err := searchServer(&allServers, s.service, v, true)
			if err != nil {
				return nil, err
			}
			s.params.UUID = server.UUID
			s.params.Timeout = timeout
			restartServers = append(restartServers, s.params.RestartServerRequest)
		}
		var (
			mu            sync.Mutex
			numOk         int
			serverDetails []*upcloud.ServerDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			req := restartServers[idx]
			msg := fmt.Sprintf("Restarting %q", req.UUID)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.RestartServer(&req)

			if err == nil {
				e.SetMessage(fmt.Sprintf("%s: restart request sent", msg))
				_, err = WaitForServerState(s.service, req.UUID, upcloud.ServerStateStarted, s.Config().ClientTimeout())
			}

			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				e.SetMessage(fmt.Sprintf("%s: done", msg))
				e.SetDetails(details.UUID, "UUID: ")
				mu.Lock()
				numOk++
				serverDetails = append(serverDetails, details)
				mu.Unlock()
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(restartServers),
			MaxConcurrentTasks: maxServerActions,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)

		if numOk != len(restartServers) {
			return nil, fmt.Errorf("number of servers failed to start: %d", len(restartServers)-int(numOk))
		}
		return serverDetails, nil
	}
}
