package objectstorage

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// ShowCommand creates the "objectstorage show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show Managed object storage service details",
			"upctl objectstorage show 55199a44-4751-4e27-9394-7c7661910be8",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	objectStorage, err := svc.GetManagedObjectStorage(exec.Context(), &request.GetManagedObjectStorageRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	endpointRows := make([]output.TableRow, 0)
	for _, endpoint := range objectStorage.Endpoints {
		endpointRows = append(endpointRows, output.TableRow{
			endpoint.DomainName,
			endpoint.Type,
		})
	}

	endpointColumns := []output.TableColumn{
		{Key: "domain_name", Header: "Domain"},
		{Key: "type", Header: "Type"},
	}

	networkRows := make([]output.TableRow, 0)
	for _, network := range objectStorage.Networks {
		networkUUID := ""
		if network.UUID != nil {
			networkUUID = *network.UUID
		}
		networkRows = append(networkRows, output.TableRow{
			network.Name,
			networkUUID,
			network.Type,
			network.Family,
		})
	}

	networkColumns := []output.TableColumn{
		{Key: "name", Header: "Name"},
		{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours, Format: format.PossiblyUnknownString},
		{Key: "type", Header: "Type"},
		{Key: "Family", Header: "Family"},
	}

	userRows := make([]output.TableRow, 0)
	for _, user := range objectStorage.Users {
		userRows = append(userRows, output.TableRow{
			user.Username,
			user.CreatedAt,
			user.UpdatedAt,
			user.OperationalState,
		})
	}

	userColumns := []output.TableColumn{
		{Key: "name", Header: "Username"},
		{Key: "created_at", Header: "Created"},
		{Key: "updated_at", Header: "Updated"},
		{Key: "operational_state", Header: "Updated", Format: format.ObjectStorageUserOperationalState},
	}

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value: objectStorage,
		Output: output.Combined{
			output.CombinedSection{
				Contents: output.Details{
					Sections: []output.DetailSection{
						{
							Title: "Overview:",
							Rows: []output.DetailRow{
								{Title: "UUID:", Value: objectStorage.UUID, Colour: ui.DefaultUUUIDColours},
								{Title: "Region:", Value: objectStorage.Region},
								{Title: "Configured status:", Value: objectStorage.ConfiguredStatus, Format: format.ObjectStorageConfiguredStatus},
								{Title: "Operational state:", Value: objectStorage.OperationalState, Format: format.ObjectStorageOperationalState},
								{Title: "Created:", Value: objectStorage.CreatedAt},
								{Title: "Updated:", Value: objectStorage.UpdatedAt},
							},
						},
					},
				},
			},
			labels.GetLabelsSectionWithResourceType(objectStorage.Labels, "managed object storage"),
			output.CombinedSection{
				Key:   "endpoints",
				Title: "Endpoints:",
				Contents: output.Table{
					Columns:      endpointColumns,
					Rows:         endpointRows,
					EmptyMessage: "No endpoints found for this Managed object storage service.",
				},
			},
			output.CombinedSection{
				Key:   "networks",
				Title: "Networks:",
				Contents: output.Table{
					Columns:      networkColumns,
					Rows:         networkRows,
					EmptyMessage: "No networks found for this Managed object storage service.",
				},
			},
			output.CombinedSection{
				Key:   "users",
				Title: "Users:",
				Contents: output.Table{
					Columns:      userColumns,
					Rows:         userRows,
					EmptyMessage: "No users found for this Managed object storage service.",
				},
			},
		},
	}, nil
}
