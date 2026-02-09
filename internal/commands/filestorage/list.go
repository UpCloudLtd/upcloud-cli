package filestorage

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/paging"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ListCommand creates the "objectstorage list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List file storage services", "upctl file-storage list"),
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
	filestorages, err := svc.GetFileStorages(exec.Context(), &request.GetFileStoragesRequest{
		Page: c.Page(),
	})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, filestorage := range filestorages {
		rows = append(rows, output.TableRow{
			filestorage.UUID,
			filestorage.Name,
			filestorage.Zone,
			filestorage.OperationalState,
			filestorage.SizeGiB,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: filestorages,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "name", Header: "Name"},
				{Key: "zone", Header: "Zone"},
				{Key: "operational_state", Header: "Operational state"},
				{Key: "size_gib", Header: "Size (GiB)"},
			},
			Rows: rows,
		},
	}, nil
}
