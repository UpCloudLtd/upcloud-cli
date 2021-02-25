package serverfirewall

import (
	"github.com/UpCloudLtd/cli/internal/commands"
)

const (
	maxFirewallActions  = 10
	positionalArgHelp = "/server/<UUID/firewall_rule>"
)

func BaseServerFirewallCommand() commands.Command {
	return &serverFirewallCommand{commands.New("firewall", "Manage server firewall")}
}

type serverFirewallCommand struct {
	*commands.BaseCommand
}
