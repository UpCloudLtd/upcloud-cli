package auditlog

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseAuditLogCommand creates the base "audit-log" command
func BaseAuditLogCommand() commands.Command {
	return &auditLogCommand{
		commands.New("audit-log", "Manage audit logs"),
	}
}

type auditLogCommand struct {
	*commands.BaseCommand
}

// InitCommand implements [commands.BaseCommand.InitCommand].
func (c *auditLogCommand) InitCommand() {
	c.Cobra().Aliases = []string{"auditlog", "al"}
}
