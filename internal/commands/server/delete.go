package server

import (
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
	deleteStorages string
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
	flags.StringVar(&s.deleteStorages, "delete-storages", "true", "Delete storages that are attached to the server.")
	s.AddFlags(flags)
	s.SetPositionalArgHelp("<uuidHostnameOrTitle ...>")
}

func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		var action = func(req interface{}) (interface{}, error) {
			server := req.(*upcloud.Server)
			var err error
			if s.deleteStorages == "true" {
				err = s.service.DeleteServerAndStorages(&request.DeleteServerAndStoragesRequest{
					UUID: server.UUID,
				})
			} else {
				err = s.service.DeleteServer(&request.DeleteServerRequest{
					UUID: server.UUID,
				})
			}
			return nil, err
		}

		return Request{
			BuildRequest: func(server *upcloud.Server) interface{} {return server},
			Service:    s.service,
			HandleContext: ui.HandleContext{
				RequestId:     func(in interface{}) string { return in.(*upcloud.Server).UUID },
				InteractiveUi: s.Config().InteractiveUI(),
				MaxActions:    maxServerActions,
				ActionMsg:     "Deleting",
				Action:        action,
			},
		}.Send(args)
	}
}
