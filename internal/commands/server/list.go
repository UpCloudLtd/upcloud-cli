package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/table"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

func ListCommand() commands.Command {
	return &listCommand{
		Command: commands.New("list", "List current servers"),
	}
}

type listCommand struct {
	commands.Command
	service        *service.Service
	headerNames    []string
	columnKeys     []string
	visibleColumns []string
}

func (s *listCommand) InitCommand() {
	s.headerNames = []string{"UUID", "Hostname", "Plan", "Zone", "State", "Tags", "Title", "License"}
	s.columnKeys = []string{"uuid", "hostname", "plan", "zone", "state", "tags", "title", "license"}
	s.visibleColumns = []string{"uuid", "hostname", "plan", "zone", "state"}
	flags := &pflag.FlagSet{}
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

func (s *listCommand) MakeExecuteCommand() func(args []string) error {
	return func(args []string) error {
		service := upapi.Service(s.Config())
		servers, err := service.GetServers()
		if err != nil {
			return err
		}
		return s.HandleOutput(servers)
	}
}

func (s *listCommand) HandleOutput(out interface{}) error {
	if s.FQConfigValueString("output") != "human" {
		return s.Command.HandleOutput(out)
	}
	servers := out.(*upcloud.Servers)
	fmt.Println()
	t, err := table.New(os.Stdout, s.headerNames...)
	if err != nil {
		return err
	}
	t.SetVisibleColumns(s.visibleColumns...)
	t.SetColumnKeys(s.columnKeys...)

	for _, server := range servers.Servers {
		t.NextRow()
		plan := server.Plan
		if plan == "custom" {
			memory := server.MemoryAmount / 1024
			plan = fmt.Sprintf("Custom (%dxCPU, %dGB)", server.CoreNumber, memory)
		}
		for _, val := range []string{server.UUID, server.Hostname, plan, server.Zone} {
			t.AddColumn(val, nil)
		}
		switch server.State {
		case upcloud.ServerStateStarted:
			t.AddColumn(server.State, tablewriter.Colors{tablewriter.FgGreenColor})
		case upcloud.ServerStateError:
			t.AddColumn(server.State, tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold})
		case upcloud.ServerStateMaintenance:
			t.AddColumn(server.State, tablewriter.Colors{tablewriter.FgYellowColor})
		default:
			t.AddColumn(server.State, tablewriter.Colors{tablewriter.FgHiBlackColor})
		}
		for _, val := range []string{strings.Join(server.Tags, ","), server.Title, fmt.Sprintf("%f", server.License)} {
			t.AddColumn(val, nil)
		}
	}
	t.Render()
	fmt.Println()
	return nil
}
