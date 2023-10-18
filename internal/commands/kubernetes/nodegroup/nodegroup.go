package nodegroup

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseNodeGroupCommand creates the base "kubernetes nodegroups" command
func BaseNodeGroupCommand() commands.Command {
	return &nodegroupCommand{
		commands.New("nodegroup", "Manage cluster node-groups"),
	}
}

type nodegroupCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (k *nodegroupCommand) InitCommand() {
	k.Cobra().Aliases = []string{"node-group"}
}
