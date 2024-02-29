package storage

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

// ListCommand creates the "storage list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List current storages",
			"upctl storage list",
			"upctl storage list --all",
		),
	}
}

type listCommand struct {
	*commands.BaseCommand
	all      config.OptionalBoolean
	private  config.OptionalBoolean
	public   config.OptionalBoolean
	normal   config.OptionalBoolean
	backup   config.OptionalBoolean
	cdrom    config.OptionalBoolean
	template config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &s.all, "all", false, "Show all storages.")
	config.AddToggleFlag(flags, &s.private, "private", true, "Show private storages (default).")
	config.AddToggleFlag(flags, &s.public, "public", false, "Show public storages.")
	config.AddToggleFlag(flags, &s.normal, "normal", false, "Show only normal storages.")
	config.AddToggleFlag(flags, &s.backup, "backup", false, "Show only backup storages.")
	config.AddToggleFlag(flags, &s.cdrom, "cdrom", false, "Show only cdrom storages.")
	config.AddToggleFlag(flags, &s.template, "template", false, "Show only template storages.")

	s.AddFlags(flags)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.Storage()

	storageList, err := svc.GetStorages(exec.Context(), &request.GetStoragesRequest{})
	if err != nil {
		return nil, err
	}

	CachedStorages = storageList.Storages
	filtered := make([]upcloud.Storage, 0)
	for _, v := range storageList.Storages {
		if s.all.Value() {
			filtered = append(filtered, v)
			continue
		}

		if s.public.Value() {
			s.private = config.False
		}

		if s.private.Value() && v.Access == upcloud.StorageAccessPublic {
			continue
		}
		if s.public.Value() && v.Access == upcloud.StorageAccessPrivate {
			continue
		}
		if !s.normal.Value() && !s.backup.Value() && !s.cdrom.Value() && !s.template.Value() {
			filtered = append(filtered, v)
			continue
		}
		if s.normal.Value() && v.Type == upcloud.StorageTypeNormal {
			filtered = append(filtered, v)
		}
		if s.backup.Value() && v.Type == upcloud.StorageTypeBackup {
			filtered = append(filtered, v)
		}
		if s.cdrom.Value() && v.Type == upcloud.StorageTypeCDROM {
			filtered = append(filtered, v)
		}
		if s.template.Value() && v.Type == upcloud.StorageTypeTemplate {
			filtered = append(filtered, v)
		}
	}

	rows := []output.TableRow{}
	for _, storage := range filtered {
		rows = append(rows, output.TableRow{
			storage.UUID,
			storage.Title,
			storage.Encrypted,
			storage.Type,
			storage.Size,
			storage.State,
			storage.Tier,
			storage.Zone,
			storage.Access,
			storage.Created,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: upcloud.Storages{Storages: filtered},
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "title", Header: "Title"},
				{Key: "encrypted", Header: "Encrypted", Format: format.Boolean},
				{Key: "type", Header: "Type"},
				{Key: "size", Header: "Size"},
				{Key: "state", Header: "State", Format: format.StorageState},
				{Key: "tier", Header: "Tier"},
				{Key: "zone", Header: "Zone"},
				{Key: "access", Header: "Access"},
			},
			Rows: rows,
		},
	}, nil
}
