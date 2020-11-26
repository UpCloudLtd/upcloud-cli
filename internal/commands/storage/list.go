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
	public         bool
	private        bool
	normal         bool
	backup         bool
	cdrom          bool
	template       bool
	favorite       bool
}

func (s *listCommand) InitCommand() {
	s.header = table.Row{"UUID", "Title", "Zone", "State", "Type", "Size", "Tier", "Created", "Access"}
	s.columnKeys = []string{"uuid", "title", "zone", "state", "type", "size", "tier", "created", "access"}
	s.visibleColumns = []string{"uuid", "title", "zone", "state", "type", "size", "tier", "created"}
	flags := &pflag.FlagSet{}
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	flags.BoolVar(&s.public, "public", false, "List public storages")
	flags.BoolVar(&s.private, "private", false, "List private storages")

	flags.BoolVar(&s.normal, "normal", false, "Filters for normal storages")
	flags.BoolVar(&s.backup, "backup", false, "Filters for backup storages")
	flags.BoolVar(&s.cdrom, "cdrom", false, "Filters for cdrom storages")
	flags.BoolVar(&s.template, "template", false, "Filters for template storages")

	s.AddFlags(flags)
}

func (s *listCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		storages, err := s.service.GetStorages(&request.GetStoragesRequest{})
		cachedStorages = storages.Storages
		if err != nil {
			return nil, err
		}
		var filtered []upcloud.Storage
		for _, v := range storages.Storages {
			if !s.public && !s.private {
				s.private = true
			}
			if s.private && !s.public && v.Access == upcloud.StorageAccessPublic {
				continue
			}
			if s.public && !s.private && v.Access == upcloud.StorageAccessPrivate {
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

		storages.Storages = filtered
		return storages, nil
	}
}

func (s *listCommand) HandleOutput(writer io.Writer, out interface{}) error {
	storages := out.(*upcloud.Storages)

	t := ui.NewDataTable(s.columnKeys...)
	t.OverrideColumnKeys(s.visibleColumns...)
	t.SetHeader(s.header)

	t.SetColumnConfig("state", table.ColumnConfig{Transformer: func(val interface{}) string {
		return StateColour(val.(string)).Sprint(val)
	}})

	sort.SliceStable(storages.Storages, func(i, j int) bool {
		return strings.Compare(storages.Storages[i].Title, storages.Storages[j].Title) < 0
	})
	sort.SliceStable(storages.Storages, func(i, j int) bool {
		return strings.Compare(storages.Storages[i].Type, storages.Storages[j].Type) < 0
	})

	for _, storage := range storages.Storages {
		t.AppendRow(table.Row{
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

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, t.Render())
	fmt.Fprintln(writer)
	return nil
}
