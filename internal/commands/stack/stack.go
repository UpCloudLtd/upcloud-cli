package stack

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

func BaseStackCommand() commands.Command {
	baseCmd := commands.New("stack", "Manage stacks")

	return &stackCommand{
		BaseCommand: baseCmd,
	}
}

type stackCommand struct {
	*commands.BaseCommand
}
