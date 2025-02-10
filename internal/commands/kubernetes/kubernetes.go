package kubernetes

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseKubernetesCommand creates the base "kubernetes" command
func BaseKubernetesCommand() commands.Command {
	return &kubernetesCommand{
		commands.New("kubernetes", "Manage Kubernetes clusters"),
	}
}

type kubernetesCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (k *kubernetesCommand) InitCommand() {
	k.Cobra().Aliases = []string{"k8s", "uks"}
	// Deprecating k8s
	// TODO: Remove this in the future
	commands.SetDeprecationHelp(k.Cobra(), []string{"k8s"})
}
