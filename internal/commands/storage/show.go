package storage

import (
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

// ShowCommand creates the "storage show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show storage details",
			"upctl storage show 01271548-2e92-44bb-9774-d282508cc762",
			"upctl storage show 01271548-2e92-44bb-9774-d282508cc762 01c60190-3a01-4108-b1c3-2e828855ccc0",
			`upctl storage show "My Storage"`,
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingStorage
	completion.Storage
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	storageSvc := exec.Storage()
	storage, err := storageSvc.GetStorageDetails(
		&request.GetStorageDetailsRequest{UUID: uuid},
	)
	if err != nil {
		return nil, err
	}

	// Storage details
	attachedToServer := "N/A"
	if len(storage.ServerUUIDs) > 0 {
		attachedToServer = strings.Join(storage.ServerUUIDs, ", \n")
	}
	storageSection := output.CombinedSection{
		Contents: output.Details{
			Sections: []output.DetailSection{
				{
					Title: "Storage",
					Rows: []output.DetailRow{
						{Title: "UUID:", Key: "uuid", Value: storage.UUID, Colour: ui.DefaultUUUIDColours},
						{Title: "Title:", Key: "title", Value: storage.Title},
						{Title: "type:", Key: "type", Value: storage.Type},
						{Title: "State:", Key: "state", Value: storage.State, Colour: commands.StorageStateColour(storage.State)},
						{Title: "Size:", Key: "size", Value: storage.Size},
						{Title: "Tier:", Key: "tier", Value: storage.Tier},
						{Title: "Zone:", Key: "zone", Value: storage.Zone},
						{Title: "Server:", Key: "server", Value: attachedToServer},
						{Title: "Origin:", Key: "origin", Value: storage.Origin, Colour: ui.DefaultUUUIDColours},
						{Title: "Created:", Key: "created", Value: storage.Created},
						{Title: "Licence:", Key: "licence", Value: storage.License},
					},
				},
			},
		},
	}

	combined := output.Combined{
		storageSection,
	}

	// Backups
	if storage.BackupRule != nil && storage.BackupRule.Interval != "" {
		combined = append(combined, output.CombinedSection{
			Contents: output.Details{
				Sections: []output.DetailSection{
					{
						Title: "Backup Rule",
						Rows: []output.DetailRow{
							{Title: "Interval:", Key: "interval", Value: storage.BackupRule.Interval},
							{Title: "Time:", Key: "time", Value: storage.BackupRule.Time},
							{Title: "Retention:", Key: "retention", Value: storage.BackupRule.Retention},
						},
					},
				},
			},
		})
	}

	if len(storage.BackupUUIDs) > 0 {
		backupsListRows := []output.TableRow{}
		for _, b := range storage.BackupUUIDs {
			backupsListRows = append(backupsListRows, output.TableRow{b})
		}
		combined = append(combined, output.CombinedSection{
			Key:   "available_backups",
			Title: "Available Backups",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				},
				Rows: backupsListRows,
			},
		})
	}

	return combined, nil
}
