package nodegroup

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseNodeGroupCommand creates the base "kubernetes nodegroups" command
func BaseNodeGroupCommand() commands.Command {
	baseCmd := commands.New("node-group", "Manage cluster node-groups")
	baseCmd.SetDeprecatedAliases([]string{"nodegroup"})
	return &nodegroupCommand{
		BaseCommand: baseCmd,
	}
}

type nodegroupCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (k *nodegroupCommand) InitCommand() {
	k.Cobra().Aliases = []string{"nodegroup"}
	commands.SetDeprecationHelp(k.Cobra(), k.DeprecatedAliases())
}
