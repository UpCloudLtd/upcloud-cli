package kubernetes

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseKubernetesCommand creates the base "kubernetes" command
func BaseKubernetesCommand() commands.Command {
	baseCmd := commands.New("kubernetes", "Manage Kubernetes clusters")
	baseCmd.SetDeprecatedAliases([]string{"uks"})

	return &kubernetesCommand{
		BaseCommand: baseCmd,
	}
}

type kubernetesCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (k *kubernetesCommand) InitCommand() {
	k.Cobra().Aliases = []string{"k8s", "uks"}

	// TODO: Remove this in the future
	commands.SetDeprecationHelp(k.Cobra(), k.DeprecatedAliases())
}
