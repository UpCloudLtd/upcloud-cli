package storage

import (
	"github.com/UpCloudLtd/cli/internal/interfaces"
	"sort"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
)

func ListCommand(service interfaces.Storage) commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current storages"),
		service:     service,
	}
}

type listCommand struct {
	*commands.BaseCommand
	service        interfaces.Storage
	header         table.Row
	columnKeys     []string
	visibleColumns []string
	public         bool
	private        bool
}

func (s *listCommand) InitCommand() {
	s.header = table.Row{"UUID", "Title", "Zone", "State", "Type", "Size", "Tier", "Created", "Access"}
	s.columnKeys = []string{"uuid", "title", "zone", "state", "type", "size", "tier", "created", "access"}
	s.visibleColumns = []string{"uuid", "title", "zone", "state", "type", "size", "tier", "created"}
	flags := &pflag.FlagSet{}
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	flags.BoolVar(&s.public, "public", false, "List public storages")
	flags.BoolVar(&s.private, "private", true, "List private storages")
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
			if !s.public && v.Access == upcloud.StorageAccessPublic {
				continue
			}
			if !s.private && v.Access == upcloud.StorageAccessPrivate {
				continue
			}
			filtered = append(filtered, v)
		}
		storages.Storages = filtered
		return storages, nil
	}
}

func (s *listCommand) HandleOutput(out interface{}) (string, error) {
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
	return t.Render(), nil
}
