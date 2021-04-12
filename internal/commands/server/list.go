package server

import (
	"fmt"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
)

// ListCommand creates the "server list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current servers"),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// Execute implements command.Command
func (s *listCommand) Execute(exec commands.Executor, _ string) (output.Output, error) {
	svc := exec.Server()
	servers, err := svc.GetServers()
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, s := range servers.Servers {
		plan := s.Plan
		if plan == "custom" {
			memory := s.MemoryAmount / 1024
			plan = fmt.Sprintf("%dxCPU-%dGB (custom)", s.CoreNumber, memory)
		}
		rows = append(rows, output.TableRow{
			s.UUID,
			s.Hostname,
			plan,
			s.Zone,
			s.State,
		})
	}

	return output.Table{
		Columns: []output.TableColumn{
			{Key: "uuid", Header: "UUID"},
			{Key: "hostname", Header: "Hostname"},
			{Key: "plan", Header: "Plan"},
			{Key: "zone", Header: "Zone"},
			{Key: "state", Header: "State"},
		},
		Rows: rows,
	}, nil
}
