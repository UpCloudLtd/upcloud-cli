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
			"Disable firewall rules by moving them after the catch-all drop rule",
			"upctl server firewall rule disable myserver --dest-port 80",
			"upctl server firewall rule disable myserver --comment \"Dev ports\"",
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
	// Get current firewall rules
	rulesResponse, err := exec.Firewall().GetFirewallRules(exec.Context(), &request.GetFirewallRulesRequest{
		ServerUUID: arg,
	})
	if err != nil {
		return nil, err
	}

	// Find the catch-all drop rule position
	catchAllPosition := findCatchAllDropRule(rulesResponse.FirewallRules)
	if catchAllPosition == 0 {
		return nil, fmt.Errorf("no catch-all drop rule found in firewall rules")
	}

	// Find matching rules that are currently before the catch-all
	matchedIndices := findMatchingRules(rulesResponse.FirewallRules, &s.params)
	var rulesToMove []int
	for _, idx := range matchedIndices {
		if rulesResponse.FirewallRules[idx].Position < catchAllPosition {
			rulesToMove = append(rulesToMove, idx)
		}
	}

	if len(rulesToMove) == 0 {
		return nil, fmt.Errorf("no enabled firewall rules matched the specified filters (rules before catch-all at position %d)", catchAllPosition)
	}

	// Confirm if multiple rules or if confirmation required
	if len(rulesToMove) > s.params.skipConfirmation {
		var ruleDescriptions []string
		for _, idx := range rulesToMove {
			rule := &rulesResponse.FirewallRules[idx]
			desc := fmt.Sprintf("  Position %d: %s %s", rule.Position, rule.Direction, rule.Protocol)
			if rule.Comment != "" {
				desc += fmt.Sprintf(" (comment: %q)", rule.Comment)
			}
			ruleDescriptions = append(ruleDescriptions, desc)
		}

		return nil, fmt.Errorf("would disable %d firewall rules (exceeds skip-confirmation=%d). Matching rules:\n%s\n\nIncrease --skip-confirmation to proceed",
			len(rulesToMove), s.params.skipConfirmation, strings.Join(ruleDescriptions, "\n"))
	}

	// Sort by position (descending) to move rules from highest position first
	sort.Slice(rulesToMove, func(i, j int) bool {
		return rulesResponse.FirewallRules[rulesToMove[i]].Position > rulesResponse.FirewallRules[rulesToMove[j]].Position
	})

	// Move each rule to after the catch-all
	movedCount := 0
	for _, idx := range rulesToMove {
		rule := &rulesResponse.FirewallRules[idx]
		oldPosition := rule.Position
		// Move to just after catch-all
		newPosition := catchAllPosition + 1

		msg := fmt.Sprintf("Disabling rule (moving from position %d to %d)", oldPosition, newPosition)
		exec.PushProgressStarted(msg)

		// Create the modified rule at the new position
		newRule := *rule
		newRule.Position = newPosition

		// Delete the old rule first
		err := exec.Firewall().DeleteFirewallRule(exec.Context(), &request.DeleteFirewallRuleRequest{
			ServerUUID: arg,
			Position:   oldPosition,
		})
		if err != nil {
			if movedCount > 0 {
				msg = fmt.Sprintf("Successfully disabled %d rules before error: %v", movedCount, err)
			}
			return commands.HandleError(exec, msg, err)
		}

		// Create the rule at the new position (after catch-all)
		_, err = exec.Firewall().CreateFirewallRule(exec.Context(), &request.CreateFirewallRuleRequest{
			ServerUUID: arg,
			FirewallRule: upcloud.FirewallRule{
				Direction:               newRule.Direction,
				Action:                  newRule.Action,
				Family:                  newRule.Family,
				Protocol:                newRule.Protocol,
				ICMPType:                newRule.ICMPType,
				DestinationAddressStart: newRule.DestinationAddressStart,
				DestinationAddressEnd:   newRule.DestinationAddressEnd,
				DestinationPortStart:    newRule.DestinationPortStart,
				DestinationPortEnd:      newRule.DestinationPortEnd,
				SourceAddressStart:      newRule.SourceAddressStart,
				SourceAddressEnd:        newRule.SourceAddressEnd,
				SourcePortStart:         newRule.SourcePortStart,
				SourcePortEnd:           newRule.SourcePortEnd,
				Comment:                 newRule.Comment,
				Position:                newPosition,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to recreate rule at position %d: %w", newPosition, err)
		}

		exec.PushProgressSuccess(msg)
		movedCount++

		// The catch-all position shifts down by 1 after each move (since we're moving from before it)
		catchAllPosition--
	}

	return output.None{}, nil
}
