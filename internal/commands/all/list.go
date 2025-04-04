package all

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/spf13/pflag"
)

// ListCommand creates the "all list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List all resources within the current account",
			"upctl all list",
			"upctl all list --include *tf-acc-test* --exclude *persistent*",
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
			{Key: "uuid", Header: "UUID", Format: formatUUID},
			{Key: "name", Header: "Name"},
		},
		Rows: rows,
	}, nil
}

func formatUUID(val interface{}) (text.Colors, string, error) {
	str, ok := val.(string)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse %T, expected string", val)
	}

	if str == "" {
		return nil, text.FgHiBlack.Sprint("-"), nil
	}
	return ui.DefaultUUUIDColours, str, nil
}
