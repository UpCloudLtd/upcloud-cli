package server

import (
	"fmt"
	"sync/atomic"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

func StartCommand() commands.Command {
	return &startCommand{
		BaseCommand: commands.New("start", "Start a server"),
	}
}

type startCommand struct {
	*commands.BaseCommand
	service   *service.Service
	avoidHost int
	host      int
}

func (s *startCommand) initService() {
	if s.service == nil {
		s.service = upapi.Service(s.Config())
	}
}

func (s *startCommand) InitCommand() {
	s.ArgCompletion(func(toComplete string) ([]string, cobra.ShellCompDirective) {
		s.initService()
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
	flags := &pflag.FlagSet{}
	flags.IntVar(&s.avoidHost, "avoid-host", 0, "Avoid specific host when starting a server")
	flags.IntVar(&s.host, "host", 0, "Start server on a specific host. Note that this is "+
		"generally available for private clouds only")
	s.AddFlags(flags)
	s.SetPositionalArgHelp("<uuidHostnameOrTitle ...>")
}

func (s *startCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		s.initService()
		if len(args) < 1 {
			return nil, fmt.Errorf("server hostname, title or uuid is required")
		}
		var (
			allServers   []upcloud.Server
			startServers []*upcloud.Server
		)
		for _, v := range args {
			server, err := searchServer(&allServers, s.service, v, true)
			if err != nil {
				return nil, err
			}
			startServers = append(startServers, server)
		}
		var numOk int64
		handler := func(idx int, e *ui.LogEntry) {
			server := startServers[idx]
			msg := fmt.Sprintf("Starting %q", server.Title)
			e.SetMessage(msg)
			e.Start()
			_, err := s.service.StartServer(&request.StartServerRequest{
				UUID:      server.UUID,
				Timeout:   s.Config().ClientTimeout(),
				AvoidHost: s.avoidHost,
				Host:      s.host,
			})
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				atomic.AddInt64(&numOk, 1)
				e.SetMessage(fmt.Sprintf("%s: done", msg))
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(startServers),
			MaxConcurrentTasks: maxServerActions,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)

		if int(numOk) < len(startServers) {
			return nil, fmt.Errorf("number of servers failed to start: %d", len(startServers)-int(numOk))
		}
		return startServers, nil
	}
}
