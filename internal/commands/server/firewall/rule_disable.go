package serverfirewall

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/spf13/pflag"
)

type ruleDisableCommand struct {
	*commands.BaseCommand
	params ruleModifyParams
	completion.Server
	resolver.CachingServer
}

// RuleDisableCommand creates the "server firewall rule disable" command
func RuleDisableCommand() commands.Command {
	return &ruleDisableCommand{
		BaseCommand: commands.New(
			"disable",
			"Disable firewall rules by changing their action to drop",
			"upctl server firewall rule disable myserver --dest-port 80",
			"upctl server firewall rule disable myserver --comment \"Dev ports\"",
			"upctl server firewall rule disable myserver --direction out --protocol udp --dest-port 53",
			"upctl server firewall rule disable myserver --position 5",
		),
		params: ruleModifyParams{
			skipConfirmation: 1,
		},
	}
}

// InitCommand implements Command.InitCommand
func (s *ruleDisableCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	addRuleFilterFlags(flagSet, &s.params, s.Cobra())
	s.AddFlags(flagSet)
	configureRuleFilterFlagsPostAdd(s.Cobra())
}

// Execute implements commands.MultipleArgumentCommand
func (s *ruleDisableCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	return modifyFirewallRules(exec, arg, &s.params, upcloud.FirewallRuleActionDrop, "disable")
}
