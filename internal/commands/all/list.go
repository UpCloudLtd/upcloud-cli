package all

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/spf13/pflag"
)

// ListCommand creates the "all list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List all resources within the current account",
			"upctl all list",
			"upctl all list --include \"*tf-acc-test*-\"",
		),
	}
}

type listCommand struct {
	*commands.BaseCommand
	include []string
	exclude []string
}

func (c *listCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.StringArrayVarP(&c.include, "include", "i", []string{"*"}, includeHelp)
	flags.StringArrayVarP(&c.exclude, "exclude", "e", []string{}, excludeHelp)
	c.AddFlags(flags)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	resources, err := listResources(exec, c.include, c.exclude)
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, resource := range resources {
		rows = append(rows, output.TableRow{
			resource.Type,
			resource.UUID,
			resource.Name,
		})
	}

	return output.Table{
		Columns: []output.TableColumn{
			{Key: "type", Header: "Type"},
			{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
			{Key: "name", Header: "Name"},
		},
		Rows: rows,
	}, nil
}
