package serverfirewall

import (
	"github.com/UpCloudLtd/cli/internal/commands"
)

const (
	PositionalArgHelp = "<server/UUID/firewall_rule>"
)

func BaseServerFirewallCommand() commands.Command {
	return &serverFirewallCommand{commands.New("firewall", "Manage server firewall rules. Enabling or disabling the firewall is done in server modify.")}
}

type serverFirewallCommand struct {
	*commands.BaseCommand
}
