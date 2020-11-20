package mocks

import "github.com/UpCloudLtd/cli/internal/commands"

func SetFlags(c commands.Command, flags ...string) error {
	return c.Cobra().Flags().Parse(flags)
}
