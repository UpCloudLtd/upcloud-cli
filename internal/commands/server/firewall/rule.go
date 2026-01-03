package serverfirewall

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseServerFirewallRuleCommand is the root command for all 'server firewall rule' commands
func BaseServerFirewallRuleCommand() commands.Command {
	return &serverFirewallRuleCommand{
		commands.New(
			"rule",
			"Manage individual firewall rules. Enable or disable specific rules by position.",
		),
	}
}

type serverFirewallRuleCommand struct {
	*commands.BaseCommand
}
