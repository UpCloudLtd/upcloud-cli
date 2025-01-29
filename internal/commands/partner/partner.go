package partner

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BasePartnerCommand creates the base "partner" command
func BasePartnerCommand() commands.Command {
	return &partnerCommand{
		commands.New("partner", "Manage partner resources"),
	}
}

type partnerCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (pr *partnerCommand) InitCommand() {
	pr.Cobra().Aliases = []string{"pr"}
}
