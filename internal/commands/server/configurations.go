package server

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"io"
)

type configurationsCommand struct {
	*commands.BaseCommand
	service service.Server
}

func ConfigurationsCommand(service service.Server) commands.Command {
	return &configurationsCommand{
		BaseCommand: commands.New("configurations", "Lists available server configurations"),
		service:     service,
	}
}

func (s *configurationsCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		configurations, err := s.service.GetServerConfigurations()
		if err != nil {
			return nil, err
		}
		return configurations, nil
	}
}

func (s *configurationsCommand) HandleOutput(writer io.Writer, out interface{}) error {
	configurations := out.(*upcloud.ServerConfigurations)

	fmt.Fprintln(writer)
	for _, cfg := range configurations.ServerConfigurations {
		fmt.Fprintln(writer, cfg.CoreNumber, ", ", cfg.MemoryAmount)
	}
	fmt.Fprintln(writer)

	return nil
}
