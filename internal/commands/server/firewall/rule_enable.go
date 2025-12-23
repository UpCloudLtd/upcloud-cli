package serverfirewall

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/spf13/pflag"
)

type ruleEnableCommand struct {
	*commands.BaseCommand
	params ruleModifyParams
	completion.Server
	resolver.CachingServer
}

// RuleEnableCommand creates the "server firewall rule enable" command
func RuleEnableCommand() commands.Command {
	return &ruleEnableCommand{
		BaseCommand: commands.New(
			"enable",
			"Enable firewall rules by changing their action to accept",
			"upctl server firewall rule enable myserver --comment \"SSH server\"",
			"upctl server firewall rule enable myserver --direction in --protocol tcp --dest-port 443",
			"upctl server firewall rule enable myserver --comment \"Dev\" --direction in --skip-confirmation 10",
			"upctl server firewall rule enable myserver --position 5",
		),
		params: ruleModifyParams{
			skipConfirmation: 1,
		},
	}
}

// InitCommand implements Command.InitCommand
func (s *ruleEnableCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	addRuleFilterFlags(flagSet, &s.params, s.Cobra())
	s.AddFlags(flagSet)
	configureRuleFilterFlagsPostAdd(s.Cobra())
}

// Execute implements commands.MultipleArgumentCommand
func (s *ruleEnableCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	return modifyFirewallRules(exec, arg, &s.params, upcloud.FirewallRuleActionAccept, "enable")
}
