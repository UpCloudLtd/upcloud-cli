package serverfirewall

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type ruleEnableCommand struct {
	*commands.BaseCommand
	rulePosition      int
	ruleComment       string
	ruleDirection     string
	ruleProtocol      string
	ruleDestPort      string
	ruleSrcAddress    string
	skipConfirmation  int
	completion.Server
	resolver.CachingServer
}

// RuleEnableCommand creates the "server firewall rule enable" command
func RuleEnableCommand() commands.Command {
	return &ruleEnableCommand{
		BaseCommand: commands.New(
			"enable",
			"Enable firewall rules by changing their action to accept",
			"upctl server firewall rule enable 00038afc-d526-4148-af0e-d2f1eeaded9b --comment \"SSH access\"",
			"upctl server firewall rule enable 00038afc-d526-4148-af0e-d2f1eeaded9b --direction in --protocol tcp --dest-port 443",
			"upctl server firewall rule enable 00038afc-d526-4148-af0e-d2f1eeaded9b --comment \"Dev\" --direction in --skip-confirmation 10",
			"upctl server firewall rule enable 00038afc-d526-4148-af0e-d2f1eeaded9b --position 5",
		),
		skipConfirmation: 1,
	}
}

// InitCommand implements Command.InitCommand
func (s *ruleEnableCommand) InitCommand() {
	directions := []string{upcloud.FirewallRuleDirectionIn, upcloud.FirewallRuleDirectionOut}
	protocols := []string{upcloud.FirewallRuleProtocolTCP, upcloud.FirewallRuleProtocolUDP, upcloud.FirewallRuleProtocolICMP}

	flagSet := &pflag.FlagSet{}
	flagSet.IntVar(&s.rulePosition, "position", 0, "Rule position. Available: 1-1000")
	flagSet.StringVar(&s.ruleComment, "comment", "", "Filter by comment (partial match, case-insensitive)")
	flagSet.StringVar(&s.ruleDirection, "direction", "", "Filter by direction. Available: "+strings.Join(directions, ", "))
	flagSet.StringVar(&s.ruleProtocol, "protocol", "", "Filter by protocol. Available: "+strings.Join(protocols, ", "))
	flagSet.StringVar(&s.ruleDestPort, "dest-port", "", "Filter by destination port (matches both start and end)")
	flagSet.StringVar(&s.ruleSrcAddress, "src-address", "", "Filter by source address (partial match)")
	flagSet.IntVar(&s.skipConfirmation, "skip-confirmation", 1, "Maximum rules to modify without confirmation. Use 0 to always require confirmation, even for a single rule.")
	s.AddFlags(flagSet)

	s.Cobra().MarkFlagsMutuallyExclusive("position", "comment")
	s.Cobra().MarkFlagsMutuallyExclusive("position", "direction")
	s.Cobra().MarkFlagsMutuallyExclusive("position", "protocol")
	s.Cobra().MarkFlagsMutuallyExclusive("position", "dest-port")
	s.Cobra().MarkFlagsMutuallyExclusive("position", "src-address")
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("position", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("comment", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("direction", cobra.FixedCompletions(directions, cobra.ShellCompDirectiveNoFileComp)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("protocol", cobra.FixedCompletions(protocols, cobra.ShellCompDirectiveNoFileComp)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("dest-port", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("src-address", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("skip-confirmation", cobra.NoFileCompletions))
}

// Execute implements commands.MultipleArgumentCommand
func (s *ruleEnableCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// Validation
	hasFilters := s.rulePosition != 0 || s.ruleComment != "" || s.ruleDirection != "" ||
		s.ruleProtocol != "" || s.ruleDestPort != "" || s.ruleSrcAddress != ""
	if !hasFilters {
		return nil, fmt.Errorf("at least one filter must be specified (--comment, --direction, --protocol, --dest-port, --src-address, or --position)")
	}
	if s.rulePosition != 0 && (s.rulePosition < 1 || s.rulePosition > 1000) {
		return nil, fmt.Errorf("invalid position (1-1000 allowed)")
	}

	msg := fmt.Sprintf("Enabling firewall rules on server %v", arg)
	exec.PushProgressStarted(msg)

	// Fetch current firewall rules
	currentRules, err := exec.Firewall().GetFirewallRules(exec.Context(), &request.GetFirewallRulesRequest{
		ServerUUID: arg,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	// Find matching rules
	var matchedIndices []int
	for i := range currentRules.FirewallRules {
		rule := &currentRules.FirewallRules[i]

		// Position-based filter (exact match, exclusive)
		if s.rulePosition != 0 {
			if rule.Position == s.rulePosition {
				matchedIndices = append(matchedIndices, i)
			}
			continue
		}

		// Apply all specified filters (AND logic)
		match := true

		if s.ruleComment != "" {
			if !strings.Contains(strings.ToLower(rule.Comment), strings.ToLower(s.ruleComment)) {
				match = false
			}
		}

		if s.ruleDirection != "" {
			if !strings.EqualFold(rule.Direction, s.ruleDirection) {
				match = false
			}
		}

		if s.ruleProtocol != "" {
			if !strings.EqualFold(rule.Protocol, s.ruleProtocol) {
				match = false
			}
		}

		if s.ruleDestPort != "" {
			// Match if either start or end matches the specified port
			if rule.DestinationPortStart != s.ruleDestPort && rule.DestinationPortEnd != s.ruleDestPort {
				match = false
			}
		}

		if s.ruleSrcAddress != "" {
			// Partial match on either start or end address
			addrLower := strings.ToLower(s.ruleSrcAddress)
			if !strings.Contains(strings.ToLower(rule.SourceAddressStart), addrLower) &&
				!strings.Contains(strings.ToLower(rule.SourceAddressEnd), addrLower) {
				match = false
			}
		}

		if match {
			matchedIndices = append(matchedIndices, i)
		}
	}

	if len(matchedIndices) == 0 {
		if s.rulePosition != 0 {
			return nil, fmt.Errorf("firewall rule at position %d not found on server %s", s.rulePosition, arg)
		}
		return nil, fmt.Errorf("no firewall rules matching the specified filters found on server %s", arg)
	}

	// Confirmation check
	if len(matchedIndices) > s.skipConfirmation {
		var ruleDescriptions []string
		for _, idx := range matchedIndices {
			rule := currentRules.FirewallRules[idx]
			desc := fmt.Sprintf("  - Position %d: %s %s", rule.Position, rule.Direction, rule.Protocol)
			if rule.Comment != "" {
				desc += fmt.Sprintf(" (comment: %q)", rule.Comment)
			}
			ruleDescriptions = append(ruleDescriptions, desc)
		}

		return nil, fmt.Errorf("would enable %d rules (exceeds skip-confirmation=%d). Matching rules:\n%s\n\nIncrease --skip-confirmation to proceed",
			len(matchedIndices), s.skipConfirmation, strings.Join(ruleDescriptions, "\n"))
	}

	// Modify matched rules
	modifiedCount := 0
	for _, idx := range matchedIndices {
		if currentRules.FirewallRules[idx].Action != upcloud.FirewallRuleActionAccept {
			currentRules.FirewallRules[idx].Action = upcloud.FirewallRuleActionAccept
			modifiedCount++
		}
	}

	if modifiedCount == 0 {
		return nil, fmt.Errorf("all %d matching rules already enabled", len(matchedIndices))
	}

	// Replace entire ruleset atomically
	err = exec.Firewall().CreateFirewallRules(exec.Context(), &request.CreateFirewallRulesRequest{
		ServerUUID:    arg,
		FirewallRules: currentRules.FirewallRules,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	msg = fmt.Sprintf("Enabled %d firewall rule(s) on server %v", modifiedCount, arg)
	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
