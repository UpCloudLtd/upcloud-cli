package all

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/spf13/pflag"
)

// PurgeCommand creates the "all purge" command
func PurgeCommand() commands.Command {
	return &purgeCommand{
		BaseCommand: commands.New(
			"purge",
			"Delete all UpCloud resources within the current account",
			"upctl all purge",
			"upctl all purge --name-prefix \"terraform-test-\"",
		),
	}
}

type purgeCommand struct {
	*commands.BaseCommand
	namePrefix string
}

func (c *purgeCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.StringVar(&c.namePrefix, "name-prefix", "", "Only delete resources having the given name prefix.")
	c.AddFlags(flags)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *purgeCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	msg := "Getting a list of all UpCloud resources"
	if c.namePrefix != "" {
		msg = fmt.Sprintf("%s having name prefix \"%s\"", msg, c.namePrefix)
	}
	exec.PushProgressStarted(msg)
	time.Sleep(2 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Deleting Managed Kubernetes clusters"
	exec.PushProgressStarted(msg)
	time.Sleep(4 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Deleting Managed Databases"
	exec.PushProgressStarted(msg)
	time.Sleep(1 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Deleting Managed Load Balancers"
	exec.PushProgressStarted(msg)
	time.Sleep(1 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Deleting Managed Object Storages"
	exec.PushProgressStarted(msg)
	time.Sleep(1 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Deleting Object Storages (deprecated)"
	exec.PushProgressStarted(msg)
	time.Sleep(4 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Deleting servers"
	exec.PushProgressStarted(msg)
	time.Sleep(4 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Deleting server groups"
	exec.PushProgressStarted(msg)
	time.Sleep(1 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Releasing IP addresses"
	exec.PushProgressStarted(msg)
	time.Sleep(1 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Deleting gateways"
	exec.PushProgressStarted(msg)
	time.Sleep(2 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Deleting routers"
	exec.PushProgressStarted(msg)
	time.Sleep(1 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Deleting networks"
	exec.PushProgressStarted(msg)
	time.Sleep(2 * time.Second)
	exec.PushProgressSuccess(msg)

	msg = "Deleting storages"
	exec.PushProgressStarted(msg)
	time.Sleep(3 * time.Second)
	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
