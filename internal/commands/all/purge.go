package all

import (
	"fmt"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/spf13/pflag"
)

// PurgeCommand creates the "all purge" command
func PurgeCommand() commands.Command {
	return &purgeCommand{
		BaseCommand: commands.New(
			"purge",
			"Delete UpCloud resources within the current account",
			"upctl all purge",
			"upctl all purge --name \"*tf-acc-test*-\"",
		),
	}
}

type purgeCommand struct {
	*commands.BaseCommand
	name string
}

func (c *purgeCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.StringVar(&c.name, "name", "", "Only delete resources matching the given name.")
	c.AddFlags(flags)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *purgeCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	msg := "Deleting all resources"
	if c.name != "" {
		msg += fmt.Sprintf(` matching name "%s"`, c.name)
	}
	exec.PushProgressStarted(msg)

	resources, err := listResources(exec, c.name)
	if err != nil {
		return nil, err
	}

	if len(resources) == 0 {
		exec.PushProgressUpdate(messages.Update{
			Key:     msg,
			Details: "Found no resources to delete",
			Status:  messages.MessageStatusWarning,
		})
		return output.None{}, nil
	}

	err = deleteResources(exec, resources, 16)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)
	return output.None{}, err
}
