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
)

func DeleteCommand(service service.Server) commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a server"),
		service:     service,
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	service        service.Server
	deleteStorages bool
}

func (s *deleteCommand) InitCommand() {
	s.ArgCompletion(func(toComplete string) ([]string, cobra.ShellCompDirective) {
		servers, err := s.service.GetServers()
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var vals []string
		for _, v := range servers.Servers {
			vals = append(vals, v.UUID, v.Hostname)
		}
		return commands.MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
	})
	flags := &pflag.FlagSet{}
	flags.BoolVar(&s.deleteStorages, "delete-storages", true, "Delete storages that are "+
		"attached to the server.")
	s.AddFlags(flags)
	s.SetPositionalArgHelp("<uuidHostnameOrTitle ...>")
}

func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("server hostname, title or uuid is required")
		}
		var (
			allServers    []upcloud.Server
			deleteServers []*upcloud.Server
		)
		for _, v := range args {
			server, err := searchServer(&allServers, s.service, v, false)
			if err != nil {
				return nil, err
			}
			deleteServers = append(deleteServers, server)
		}
		var numOk int64
		handler := func(idx int, e *ui.LogEntry) {
			server := deleteServers[idx]
			msg := fmt.Sprintf("Deleting %q", server.Title)
			e.SetMessage(msg)
			e.Start()
			var err error
			if s.deleteStorages {
				err = s.service.DeleteServerAndStorages(&request.DeleteServerAndStoragesRequest{
					UUID: server.UUID,
				})
			} else {
				err = s.service.DeleteServer(&request.DeleteServerRequest{
					UUID: server.UUID,
				})
			}
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				atomic.AddInt64(&numOk, 1)
				e.SetMessage(fmt.Sprintf("%s: done", msg))
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(deleteServers),
			MaxConcurrentTasks: maxServerActions,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)

		if int(numOk) < len(deleteServers) {
			return nil, fmt.Errorf("number of servers failed to delete: %d", len(deleteServers)-int(numOk))
		}
		return deleteServers, nil
	}
}