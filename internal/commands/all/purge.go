package all

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/spf13/pflag"
)

// PurgeCommand creates the "all purge" command
func PurgeCommand() commands.Command {
	return &purgeCommand{
		BaseCommand: commands.New(
			"purge",
			"Delete all resources from the current account",
			"upctl all purge",
			"upctl all purge --include *tf-acc-test* --exclude *persistent*",
		),
	}
}

type purgeCommand struct {
	*commands.BaseCommand
	include []string
	exclude []string
}

func (c *purgeCommand) InitCommand() {
	c.Cobra().Long = `Delete all resources from the current account. Use ` + "`" + `upctl all list` + "`" + ` command to preview targeted resources before purging.`

	flags := &pflag.FlagSet{}
	flags.StringArrayVarP(&c.include, "include", "i", []string{"*"}, includeHelp)
	flags.StringArrayVarP(&c.exclude, "exclude", "e", []string{}, excludeHelp)
	c.AddFlags(flags)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *purgeCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	resources, err := listResources(exec, c.include, c.exclude)
	if err != nil {
		return nil, err
	}

	err = deleteResources(exec, resources, 16)
	if err != nil {
		return nil, err
	}

	return output.None{}, err
}
