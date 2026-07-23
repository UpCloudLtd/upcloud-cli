package stack

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

const (
	LabelKeyStack     = "stacks.upcloud.com/stack"
	LabelKeyCreatedBy = "stacks.upcloud.com/created-by"
	LabelKeyVersion   = "stacks.upcloud.com/version"
	LabelKeyName      = "stacks.upcloud.com/name"
	LabelValueUpctl   = "upctl"
)

func BaseStackCommand() commands.Command {
	baseCmd := commands.New("stack", "Manage stacks (EXPERIMENTAL)")

	return &stackCommand{
		BaseCommand: baseCmd,
	}
}

type stackCommand struct {
	*commands.BaseCommand
}
