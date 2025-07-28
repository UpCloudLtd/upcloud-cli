package devices

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

// DevicesCommand creates the "zone devices" command
func DevicesCommand() commands.Command {
	return &devicesCommand{
		BaseCommand: commands.New("devices", "List available devices for each zone", "upctl zone devices"),
	}
}

type devicesCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *devicesCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	d, err := svc.GetDevicesAvailability(exec.Context())
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for zone, devices := range *d {
		features := []string{}
		if len(devices.GPUPlans) > 0 {
			features = append(features, "GPU")
		}

		rows = append(rows, output.TableRow{
			zone,
			features,
		})
	}

	columns := []output.TableColumn{
		{Key: "zone", Header: "Zone"},
		{Key: "devices", Header: "Devices", Format: format.StringSliceSingleLineAnd},
	}

	return output.MarshaledWithHumanOutput{
		Value: d,
		Output: output.Table{
			Columns: columns,
			Rows:    rows,
		},
	}, nil
}
