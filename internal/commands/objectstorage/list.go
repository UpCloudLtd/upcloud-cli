package objectstorage

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/paging"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ListCommand creates the "objectstorage list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current Managed object storage services", "upctl object-storage list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
	paging.PageParameters
}

func (c *listCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	c.ConfigureFlags(fs)
	c.AddFlags(fs)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	objectstorages, err := svc.GetManagedObjectStorages(exec.Context(), &request.GetManagedObjectStoragesRequest{
		Page: c.Page(),
	})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, objectstorage := range objectstorages {
		rows = append(rows, output.TableRow{
			objectstorage.UUID,
			objectstorage.Name,
			objectstorage.Region,
			objectstorage.ConfiguredStatus,
			objectstorage.OperationalState,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: objectstorages,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "name", Header: "Name"},
				{Key: "region", Header: "Region"},
				{Key: "configured_status", Header: "Configured status", Format: format.ObjectStorageConfiguredStatus},
				{Key: "operational_state", Header: "Operational state", Format: format.ObjectStorageOperationalState},
			},
			Rows: rows,
		},
	}, nil
}
