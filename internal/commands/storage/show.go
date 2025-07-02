package storage

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
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
	storage, err := storageSvc.GetStorageDetails(exec.Context(), &request.GetStorageDetailsRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	storageSection := output.CombinedSection{
		Contents: output.Details{
			Sections: []output.DetailSection{
				{
					Title: "Storage",
					Rows: []output.DetailRow{
						{Title: "UUID:", Key: "uuid", Value: storage.UUID, Colour: ui.DefaultUUUIDColours},
						{Title: "Title:", Key: "title", Value: storage.Title},
						{Title: "Access:", Key: "access", Value: storage.Access},
						{Title: "Type:", Key: "type", Value: storage.Type, Format: formatShowType(storage.TemplateType)},
						{Title: "State:", Key: "state", Value: storage.State, Format: format.StorageState},
						{Title: "Size:", Key: "size", Value: storage.Size},
						{Title: "Tier:", Key: "tier", Value: storage.Tier},
						{Title: "Encrypted:", Key: "encrypted", Value: storage.Encrypted, Format: format.Boolean},
						{Title: "Zone:", Key: "zone", Value: storage.Zone},
						{Title: "Server:", Key: "servers", Value: storage.ServerUUIDs, Format: formatShowServers},
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

	combined = append(combined, labels.GetLabelsSection(storage.Labels))

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

	return output.MarshaledWithHumanOutput{
		Value:  storage,
		Output: combined,
	}, nil
}

func formatShowServers(val interface{}) (text.Colors, string, error) {
	servers, ok := val.(upcloud.ServerUUIDSlice)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse server UUIDs from %T, expected upcloud.ServerUUIDSlice", val)
	}

	var strs []string
	for _, server := range servers {
		strs = append(strs, ui.DefaultUUUIDColours.Sprint(server))
	}

	str := "None"
	if len(servers) > 0 {
		str = strings.Join(strs, ", \n")
	}

	return nil, str, nil
}

func formatShowType(templateType string) func(any) (text.Colors, string, error) {
	return func(val any) (text.Colors, string, error) {
		st, ok := val.(string)
		if !ok {
			return nil, "", fmt.Errorf("cannot render storage type from %T, expected string", val)
		}

		if templateType == "" {
			return nil, st, nil
		}
		return nil, fmt.Sprintf("%s (%s)", st, templateType), nil
	}
}
