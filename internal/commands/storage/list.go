package storage

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"io"
	"sort"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
)

// ListCommand creates the "storage list" command
func ListCommand(service service.Storage) commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current storages"),
		service:     service,
	}
}

type listCommand struct {
	*commands.BaseCommand
	service        service.Storage
	header         table.Row
	columnKeys     []string
	visibleColumns []string
	all            bool
	private        bool
	public         bool
	normal         bool
	backup         bool
	cdrom          bool
	template       bool
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	s.header = table.Row{"UUID", "Title", "Zone", "State", "Type", "Size", "Tier", "Created", "Access"}
	s.columnKeys = []string{"uuid", "title", "zone", "state", "type", "size", "tier", "created", "access"}
	s.visibleColumns = []string{"uuid", "title", "zone", "state", "type", "size", "tier", "created"}
	flags := &pflag.FlagSet{}
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	flags.BoolVar(&s.all, "all", false, "Show all storages.")
	flags.BoolVar(&s.private, "private", true, "Show private storages (default).")
	flags.BoolVar(&s.public, "public", false, "Show public storages.")

	flags.BoolVar(&s.normal, "normal", false, "Show only normal storages.")
	flags.BoolVar(&s.backup, "backup", false, "Show only backup storages.")
	flags.BoolVar(&s.cdrom, "cdrom", false, "Show only cdrom storages.")
	flags.BoolVar(&s.template, "template", false, "Show only template storages.")

	s.AddFlags(flags)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *listCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		gotStorages, err := s.service.GetStorages(&request.GetStoragesRequest{})
		if err != nil {
			return nil, err
		}
		CachedStorages = gotStorages.Storages
		var filtered []upcloud.Storage
		for _, v := range gotStorages.Storages {
			if s.all {
				filtered = append(filtered, v)
				continue
			}

			if s.public {
				s.private = false
			}

			if s.private && v.Access == upcloud.StorageAccessPublic {
				continue
			}
			if s.public && v.Access == upcloud.StorageAccessPrivate {
				continue
			}
			if !s.normal && !s.backup && !s.cdrom && !s.template {
				filtered = append(filtered, v)
				continue
			}
			if s.normal && v.Type == upcloud.StorageTypeNormal {
				filtered = append(filtered, v)
			}
			if s.backup && v.Type == upcloud.StorageTypeBackup {
				filtered = append(filtered, v)
			}
			if s.cdrom && v.Type == upcloud.StorageTypeCDROM {
				filtered = append(filtered, v)
			}
			if s.template && v.Type == upcloud.StorageTypeTemplate {
				filtered = append(filtered, v)
			}
		}

		gotStorages.Storages = filtered
		return gotStorages, nil
	}
}

// HandleOutput implements Command.HandleOutput
func (s *listCommand) HandleOutput(writer io.Writer, out interface{}) error {
	storages := out.(*upcloud.Storages)

	t := ui.NewDataTable(s.columnKeys...)
	t.OverrideColumnKeys(s.visibleColumns...)
	t.SetHeader(s.header)

	t.SetColumnConfig("state", table.ColumnConfig{Transformer: func(val interface{}) string {
		return storageStateColor(val.(string)).Sprint(val)
	}})

	sort.SliceStable(storages.Storages, func(i, j int) bool {
		return strings.Compare(storages.Storages[i].Title, storages.Storages[j].Title) < 0
	})
	sort.SliceStable(storages.Storages, func(i, j int) bool {
		return strings.Compare(storages.Storages[i].Type, storages.Storages[j].Type) < 0
	})

	for _, storage := range storages.Storages {
		t.Append(table.Row{
			storage.UUID,
			storage.Title,
			storage.Zone,
			storage.State,
			storage.Type,
			storage.Size,
			storage.Tier,
			storage.Created,
			storage.Access})
	}
	_, _ = fmt.Fprintln(writer, t.Render())
	return nil
}
