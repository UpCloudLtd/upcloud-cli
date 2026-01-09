package serverfirewall

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ruleModifyParams holds common parameters for rule modification commands
type ruleModifyParams struct {
	rulePosition     int
	ruleComment      string
	ruleDirection    string
	ruleProtocol     string
	ruleDestPort     string
	ruleSrcAddress   string
	skipConfirmation int
}

// addRuleFilterFlags adds common filter flags to a command
// This function must be called AFTER the flagSet has been added to the cobra command
func addRuleFilterFlags(flagSet *pflag.FlagSet, params *ruleModifyParams, cobraCmd *cobra.Command) {
	directions := []string{upcloud.FirewallRuleDirectionIn, upcloud.FirewallRuleDirectionOut}
	protocols := []string{upcloud.FirewallRuleProtocolTCP, upcloud.FirewallRuleProtocolUDP, upcloud.FirewallRuleProtocolICMP}

	flagSet.IntVar(&params.rulePosition, "position", 0, "Rule position. Available: 1-1000")
	flagSet.StringVar(&params.ruleComment, "comment", "", "Filter by comment (partial match, case-insensitive)")
	flagSet.StringVar(&params.ruleDirection, "direction", "", "Filter by direction. Available: "+strings.Join(directions, ", "))
	flagSet.StringVar(&params.ruleProtocol, "protocol", "", "Filter by protocol. Available: "+strings.Join(protocols, ", "))
	flagSet.StringVar(&params.ruleDestPort, "dest-port", "", "Filter by destination port (matches both start and end)")
	flagSet.StringVar(&params.ruleSrcAddress, "src-address", "", "Filter by source address (partial match)")
	flagSet.IntVar(&params.skipConfirmation, "skip-confirmation", 1, "Maximum rules to modify without confirmation. Use 0 to always require confirmation, even for a single rule.")
}

// configureRuleFilterFlagsPostAdd configures mutual exclusivity and completion after flags are added
func configureRuleFilterFlagsPostAdd(cobraCmd *cobra.Command) {
	directions := []string{upcloud.FirewallRuleDirectionIn, upcloud.FirewallRuleDirectionOut}
	protocols := []string{upcloud.FirewallRuleProtocolTCP, upcloud.FirewallRuleProtocolUDP, upcloud.FirewallRuleProtocolICMP}

	cobraCmd.MarkFlagsMutuallyExclusive("position", "comment")
	cobraCmd.MarkFlagsMutuallyExclusive("position", "direction")
	cobraCmd.MarkFlagsMutuallyExclusive("position", "protocol")
	cobraCmd.MarkFlagsMutuallyExclusive("position", "dest-port")
	cobraCmd.MarkFlagsMutuallyExclusive("position", "src-address")

	commands.Must(cobraCmd.RegisterFlagCompletionFunc("position", cobra.NoFileCompletions))
	commands.Must(cobraCmd.RegisterFlagCompletionFunc("comment", cobra.NoFileCompletions))
	commands.Must(cobraCmd.RegisterFlagCompletionFunc("direction", cobra.FixedCompletions(directions, cobra.ShellCompDirectiveNoFileComp)))
	commands.Must(cobraCmd.RegisterFlagCompletionFunc("protocol", cobra.FixedCompletions(protocols, cobra.ShellCompDirectiveNoFileComp)))
	commands.Must(cobraCmd.RegisterFlagCompletionFunc("dest-port", completeCommonPorts))
	commands.Must(cobraCmd.RegisterFlagCompletionFunc("src-address", completeIPAddress))
	commands.Must(cobraCmd.RegisterFlagCompletionFunc("skip-confirmation", completeSkipConfirmation))
}

// findMatchingRules finds rules matching the specified filters
func findMatchingRules(rules []upcloud.FirewallRule, params *ruleModifyParams) []int {
	var matchedIndices []int

	for i := range rules {
		rule := &rules[i]

		// Position-based filter (exact match, exclusive)
		if params.rulePosition != 0 {
			if rule.Position == params.rulePosition {
				matchedIndices = append(matchedIndices, i)
			}
			continue
		}

		// Apply all specified filters (AND logic)
		match := true

		if params.ruleComment != "" {
			if !strings.Contains(strings.ToLower(rule.Comment), strings.ToLower(params.ruleComment)) {
				match = false
			}
		}

		if params.ruleDirection != "" {
			if !strings.EqualFold(rule.Direction, params.ruleDirection) {
				match = false
			}
		}

		if params.ruleProtocol != "" {
			if !strings.EqualFold(rule.Protocol, params.ruleProtocol) {
				match = false
			}
		}

		if params.ruleDestPort != "" {
			// Match if either start or end matches the specified port
			if rule.DestinationPortStart != params.ruleDestPort && rule.DestinationPortEnd != params.ruleDestPort {
				match = false
			}
		}

		if params.ruleSrcAddress != "" {
			// Partial match on either start or end address
			addrLower := strings.ToLower(params.ruleSrcAddress)
			if !strings.Contains(strings.ToLower(rule.SourceAddressStart), addrLower) &&
				!strings.Contains(strings.ToLower(rule.SourceAddressEnd), addrLower) {
				match = false
			}
		}

		if match {
			matchedIndices = append(matchedIndices, i)
		}
	}

	return matchedIndices
}

// modifyFirewallRules modifies firewall rules based on filters and applies the given action
func modifyFirewallRules(
	exec commands.Executor,
	serverUUID string,
	params *ruleModifyParams,
	targetAction string,
	actionVerb string, // "enable" or "disable"
) (output.Output, error) {
	// Validation
	hasFilters := params.rulePosition != 0 || params.ruleComment != "" || params.ruleDirection != "" ||
		params.ruleProtocol != "" || params.ruleDestPort != "" || params.ruleSrcAddress != ""
	if !hasFilters {
		return nil, fmt.Errorf("at least one filter must be specified (--comment, --direction, --protocol, --dest-port, --src-address, or --position)")
	}
	if params.rulePosition != 0 && (params.rulePosition < 1 || params.rulePosition > 1000) {
		return nil, fmt.Errorf("invalid position (1-1000 allowed)")
	}

	msg := fmt.Sprintf("%sing firewall rules on server %v", strings.Title(actionVerb), serverUUID)
	exec.PushProgressStarted(msg)

	// Fetch current firewall rules
	currentRules, err := exec.Firewall().GetFirewallRules(exec.Context(), &request.GetFirewallRulesRequest{
		ServerUUID: serverUUID,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	// Find matching rules
	matchedIndices := findMatchingRules(currentRules.FirewallRules, params)

	if len(matchedIndices) == 0 {
		if params.rulePosition != 0 {
			return nil, fmt.Errorf("firewall rule at position %d not found on server %s", params.rulePosition, serverUUID)
		}
		return nil, fmt.Errorf("no firewall rules matching the specified filters found on server %s", serverUUID)
	}

	// Confirmation check
	if len(matchedIndices) > params.skipConfirmation {
		var ruleDescriptions []string
		for _, idx := range matchedIndices {
			rule := currentRules.FirewallRules[idx]
			desc := fmt.Sprintf("  - Position %d: %s %s", rule.Position, rule.Direction, rule.Protocol)
			if rule.Comment != "" {
				desc += fmt.Sprintf(" (comment: %q)", rule.Comment)
			}
			ruleDescriptions = append(ruleDescriptions, desc)
		}

		return nil, fmt.Errorf("would %s %d rules (exceeds skip-confirmation=%d). Matching rules:\n%s\n\nIncrease --skip-confirmation to proceed",
			actionVerb, len(matchedIndices), params.skipConfirmation, strings.Join(ruleDescriptions, "\n"))
	}

	// Modify matched rules
	modifiedCount := 0
	for _, idx := range matchedIndices {
		if currentRules.FirewallRules[idx].Action != targetAction {
			currentRules.FirewallRules[idx].Action = targetAction
			modifiedCount++
		}
	}

	if modifiedCount == 0 {
		return nil, fmt.Errorf("all %d matching rules already %sd", len(matchedIndices), actionVerb)
	}

	// Replace entire ruleset atomically
	err = exec.Firewall().CreateFirewallRules(exec.Context(), &request.CreateFirewallRulesRequest{
		ServerUUID:    serverUUID,
		FirewallRules: currentRules.FirewallRules,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	msg = fmt.Sprintf("%sd %d firewall rule(s) on server %v", strings.Title(actionVerb), modifiedCount, serverUUID)
	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
