package serverfirewall

import (
	"github.com/UpCloudLtd/cli/internal/commands"
)

const (
// TODO: re-add
// positionalArgHelp = "<server/UUID/firewall_rule>"
)

// BaseServerFirewallCommand is the root command for all 'server firewall' commands
func BaseServerFirewallCommand() commands.Command {
	return &serverFirewallCommand{commands.New("firewall", "Manage server firewall rules. Enabling or disabling the firewall is done in server modify.")}
}

type serverFirewallCommand struct {
	*commands.BaseCommand
}
