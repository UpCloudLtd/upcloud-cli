package objectstorage

import (
	"sort"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// RegionsCommand creates the "objectstorage regions" command
func RegionsCommand() commands.Command {
	return &regionsCommand{
		BaseCommand: commands.New("regions", "List objectstorage regions", "upctl objectstorage regions"),
	}
}

type regionsCommand struct {
	*commands.BaseCommand
}

func (s *regionsCommand) InitCommand() {
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *regionsCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	regions, err := exec.All().GetManagedObjectStorageRegions(exec.Context(), &request.GetManagedObjectStorageRegionsRequest{})
	if err != nil {
		return nil, err
	}

	sort.Slice(regions, func(i, j int) bool {
		return regions[i].Name < regions[j].Name
	})

	rows := []output.TableRow{}
	for _, r := range regions {
		zones := []string{}
		for _, z := range r.Zones {
			zones = append(zones, z.Name)
		}
		sort.Strings(zones)

		rows = append(rows, output.TableRow{
			r.Name,
			r.PrimaryZone,
			zones,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: regions,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "name", Header: "Name"},
				{Key: "primary_zone", Header: "Primary zone"},
				{Key: "zones", Header: "Zones", Format: format.StringSliceSingleLineAnd},
			},
			Rows: rows,
		},
	}, nil
}
