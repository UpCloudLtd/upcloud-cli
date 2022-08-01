package serverfirewall

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
)

// BaseServerFirewallCommand is the root command for all 'server firewall' commands
func BaseServerFirewallCommand() commands.Command {
	return &serverFirewallCommand{
		commands.New(
			"firewall",
			"Manage server firewall rules. Enabling or disabling the firewall is done in server modify.",
		)}
}

type serverFirewallCommand struct {
	*commands.BaseCommand
}

func (s *serverFirewallCommand) BuildSubCommands(cfg *config.Config) {
	commands.BuildCommand(CreateCommand(), s.Cobra(), cfg)
	commands.BuildCommand(DeleteCommand(), s.Cobra(), cfg)
	commands.BuildCommand(ShowCommand(), s.Cobra(), cfg)
}
