package server

import (
	"fmt"
	"os"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/cobra"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/upapi"
	"github.com/UpCloudLtd/cli/internal/validation"
)

func ShowCommand() commands.Command {
	return &showCommand{
		Command: commands.New("show", "Show server details"),
	}
}

type showCommand struct {
	commands.Command
	service *service.Service
}

func (s *showCommand) initService() {
	if s.service == nil {
		s.service = upapi.Service(s.Config())
	}
}

func (s *showCommand) InitCommand() {
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
		fmt.Fprintf(os.Stderr, "%v\n", vals)
		return commands.MatchStringPrefix(vals, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
}

func (s *showCommand) matchServer(servers []upcloud.Server, searchVal string) *upcloud.Server {
	for _, server := range servers {
		if server.Title == searchVal || server.Hostname == searchVal {
			return &server
		}
	}
	return nil
}

func (s *showCommand) MakeExecuteCommand() func(args []string) error {
	return func(args []string) error {
		s.initService()
		// TODO(aakso): implement prompting with readline support
		if len(args) < 1 {
			return fmt.Errorf("server hostname, title or uuid is required")
		}
		serverUuid := args[0]
		if err := validation.Uuid4(args[0]); err != nil {
			servers, err := s.service.GetServers()
			if err != nil {
				return err
			}
			server := s.matchServer(servers.Servers, args[0])
			if server == nil {
				return fmt.Errorf("no server with name or title %q was found", args[0])
			}
			serverUuid = server.UUID
		}
		server, err := s.service.GetServerDetails(&request.GetServerDetailsRequest{UUID: serverUuid})
		if err != nil {
			return err
		}
		s.HandleOutput(server)
		return nil
	}
}

func (s *showCommand) HandleOutput(out interface{}) error {
	if s.Config().GetString("output") != "human" {
		return s.Command.HandleOutput(out)
	}
	return nil
}
