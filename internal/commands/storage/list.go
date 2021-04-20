package storage

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

// ListCommand creates the "storage list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current storages", ""),
	}
}

type listCommand struct {
	*commands.BaseCommand
	all      bool
	private  bool
	public   bool
	normal   bool
	backup   bool
	cdrom    bool
	template bool
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.BoolVar(&s.all, "all", false, "Show all storages.")
	flags.BoolVar(&s.private, "private", true, "Show private storages (default).")
	flags.BoolVar(&s.public, "public", false, "Show public storages.")
	flags.BoolVar(&s.normal, "normal", false, "Show only normal storages.")
	flags.BoolVar(&s.backup, "backup", false, "Show only backup storages.")
	flags.BoolVar(&s.cdrom, "cdrom", false, "Show only cdrom storages.")
	flags.BoolVar(&s.template, "template", false, "Show only template storages.")

	s.AddFlags(flags)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.Storage()

	storageList, err := svc.GetStorages(&request.GetStoragesRequest{})
	if err != nil {
		return nil, err
	}

	CachedStorages = storageList.Storages
	var filtered []upcloud.Storage
	for _, v := range storageList.Storages {
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

	rows := []output.TableRow{}
	for _, storage := range filtered {
		rows = append(rows, output.TableRow{
			storage.UUID,
			storage.Title,
			storage.Type,
			storage.Size,
			storage.State,
			storage.Tier,
			storage.Zone,
			storage.Access,
			storage.Created,
		})
	}

	return output.Table{
		Columns: []output.TableColumn{
			{Key: "uuid", Header: "UUID", Color: ui.DefaultUUUIDColours},
			{Key: "title", Header: "Title"},
			{Key: "type", Header: "Type"},
			{Key: "size", Header: "Size"},
			{Key: "state", Header: "State"},
			{Key: "tier", Header: "Tier"},
			{Key: "zone", Header: "zone"},
			{Key: "access", Header: "Access"},
			{Key: "created", Header: "Created"},
		},
		Rows: rows,
	}, nil
}
