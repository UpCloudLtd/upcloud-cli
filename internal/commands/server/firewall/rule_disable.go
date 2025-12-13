package serverfirewall

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type ruleDisableCommand struct {
	*commands.BaseCommand
	rulePosition int
	completion.Server
	resolver.CachingServer
}

// RuleDisableCommand creates the "server firewall rule disable" command
func RuleDisableCommand() commands.Command {
	return &ruleDisableCommand{
		BaseCommand: commands.New(
			"disable",
			"Disable a specific firewall rule by changing its action to drop",
			"upctl server firewall rule disable 00038afc-d526-4148-af0e-d2f1eeaded9b --position 5",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *ruleDisableCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.IntVar(&s.rulePosition, "position", 0, "Rule position. Available: 1-1000")
	s.AddFlags(flagSet)

	commands.Must(s.Cobra().MarkFlagRequired("position"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("position", cobra.NoFileCompletions))
}

// Execute implements commands.MultipleArgumentCommand
func (s *ruleDisableCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if s.rulePosition < 1 || s.rulePosition > 1000 {
		return nil, fmt.Errorf("invalid position (1-1000 allowed)")
	}

	msg := fmt.Sprintf("Disabling firewall rule %d on server %v", s.rulePosition, arg)
	exec.PushProgressStarted(msg)

	// Fetch current firewall rules
	currentRules, err := exec.Firewall().GetFirewallRules(exec.Context(), &request.GetFirewallRulesRequest{
		ServerUUID: arg,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	// Find and modify the target rule
	ruleFound := false
	for i := range currentRules.FirewallRules {
		if currentRules.FirewallRules[i].Position == s.rulePosition {
			currentRules.FirewallRules[i].Action = upcloud.FirewallRuleActionDrop
			ruleFound = true
			break
		}
	}

	if !ruleFound {
		return nil, fmt.Errorf("firewall rule at position %d not found on server %s", s.rulePosition, arg)
	}

	// Replace entire ruleset atomically
	err = exec.Firewall().CreateFirewallRules(exec.Context(), &request.CreateFirewallRulesRequest{
		ServerUUID:    arg,
		FirewallRules: currentRules.FirewallRules,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
