package serverfirewall

import (
	"fmt"
	"sort"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

type deleteCommand struct {
	*commands.BaseCommand
	params ruleModifyParams
	completion.Server
	resolver.CachingServer
}

// DeleteCommand creates the "server firewall delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete firewall rules from a server. Rules can be deleted by position or by using filters.",
			"upctl server firewall delete 00038afc-d526-4148-af0e-d2f1eeaded9b --position 1",
			"upctl server firewall delete myserver --comment \"temporary rule\"",
			"upctl server firewall delete myserver --direction in --protocol tcp --dest-port 8080",
		),
		params: ruleModifyParams{
			skipConfirmation: 1,
		},
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	addRuleFilterFlags(flagSet, &s.params, s.Cobra())
	s.AddFlags(flagSet)
	configureRuleFilterFlagsPostAdd(s.Cobra())
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// Get current firewall rules
	rulesResponse, err := exec.Firewall().GetFirewallRules(exec.Context(), &request.GetFirewallRulesRequest{
		ServerUUID: arg,
	})
	if err != nil {
		return nil, err
	}

	// Find matching rules
	matchedIndices := findMatchingRules(rulesResponse.FirewallRules, &s.params)

	if len(matchedIndices) == 0 {
		return nil, fmt.Errorf("no firewall rules matched the specified filters")
	}

	// Confirm if multiple rules or if confirmation required
	if len(matchedIndices) > s.params.skipConfirmation {
		var ruleDescriptions []string
		for _, idx := range matchedIndices {
			rule := &rulesResponse.FirewallRules[idx]
			desc := fmt.Sprintf("  Position %d: %s %s %s -> %s",
				rule.Position, rule.Direction, rule.Protocol,
				formatRuleAddress(rule, true), formatRuleAddress(rule, false))
			if rule.Comment != "" {
				desc += fmt.Sprintf(" (comment: %q)", rule.Comment)
			}
			ruleDescriptions = append(ruleDescriptions, desc)
		}

		return nil, fmt.Errorf("would delete %d firewall rules (exceeds skip-confirmation=%d). Matching rules:\n%s\n\nIncrease --skip-confirmation to proceed",
			len(matchedIndices), s.params.skipConfirmation, strings.Join(ruleDescriptions, "\n"))
	}

	// Sort indices in descending order to delete from highest position first
	// This prevents position shifts affecting subsequent deletions
	sort.Sort(sort.Reverse(sort.IntSlice(matchedIndices)))

	// Delete each matched rule
	deletedCount := 0
	for _, idx := range matchedIndices {
		rule := &rulesResponse.FirewallRules[idx]
		msg := fmt.Sprintf("Deleting firewall rule at position %d", rule.Position)
		exec.PushProgressStarted(msg)

		err := exec.Firewall().DeleteFirewallRule(exec.Context(), &request.DeleteFirewallRuleRequest{
			ServerUUID: arg,
			Position:   rule.Position,
		})
		if err != nil {
			if deletedCount > 0 {
				msg = fmt.Sprintf("Successfully deleted %d rules before error: %v", deletedCount, err)
			}
			return commands.HandleError(exec, msg, err)
		}

		exec.PushProgressSuccess(msg)
		deletedCount++

		// Adjust positions of remaining rules in our local copy
		for i := range rulesResponse.FirewallRules {
			if rulesResponse.FirewallRules[i].Position > rule.Position {
				rulesResponse.FirewallRules[i].Position--
			}
		}
	}

	return output.None{}, nil
}

func formatRuleAddress(rule *upcloud.FirewallRule, source bool) string {
	var addr, port string
	if source {
		addr = rule.SourceAddressStart
		if rule.SourceAddressEnd != "" && rule.SourceAddressEnd != rule.SourceAddressStart {
			addr = fmt.Sprintf("%s-%s", addr, rule.SourceAddressEnd)
		}
		port = rule.SourcePortStart
		if rule.SourcePortEnd != "" && rule.SourcePortEnd != rule.SourcePortStart {
			port = fmt.Sprintf("%s-%s", port, rule.SourcePortEnd)
		}
	} else {
		addr = rule.DestinationAddressStart
		if rule.DestinationAddressEnd != "" && rule.DestinationAddressEnd != rule.DestinationAddressStart {
			addr = fmt.Sprintf("%s-%s", addr, rule.DestinationAddressEnd)
		}
		port = rule.DestinationPortStart
		if rule.DestinationPortEnd != "" && rule.DestinationPortEnd != rule.DestinationPortStart {
			port = fmt.Sprintf("%s-%s", port, rule.DestinationPortEnd)
		}
	}

	if port != "" {
		return fmt.Sprintf("%s:%s", addr, port)
	}
	return addr
}
