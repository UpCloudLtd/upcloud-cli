package devices

import (
	"sort"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

// ShowCommand creates the "zone devices show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show detailed device availability for the given zone", "upctl zone devices show fi-hel2"),
	}
}

type showCommand struct {
	*commands.BaseCommand
	completion.Zone
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, zone string) (output.Output, error) {
	svc := exec.All()
	d, err := svc.GetDevicesAvailability(exec.Context())
	if err != nil {
		return nil, err
	}

	sections := output.Combined{}

	if len((*d)[zone].GPUPlans) > 0 {
		rows := []output.TableRow{}
		for name, data := range (*d)[zone].GPUPlans {
			rows = append(rows, output.TableRow{
				name,
				data.Amount,
			})
		}

		sections = append(sections, output.CombinedSection{
			Key:   "gpu_plans",
			Title: "GPU plans",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "name", Header: "Name"},
					{Key: "amount", Header: "Amount"},
				},
				Rows: rows,
			},
		})

		sort.Slice(rows, func(i, j int) bool {
			return rows[i][0].(string) < rows[j][0].(string)
		})
	}

	return output.MarshaledWithHumanOutput{
		Value:  (*d)[zone],
		Output: sections,
	}, nil
}
