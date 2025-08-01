package serverfirewall

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type deleteCommand struct {
	*commands.BaseCommand
	rulePosition int
	completion.Server
	resolver.CachingServer
}

// DeleteCommand creates the "server firewall delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Removes a firewall rule from a server. Firewall rules must be removed individually. The positions of remaining firewall rules will be adjusted after a rule is removed.",
			"upctl server firewall delete 00038afc-d526-4148-af0e-d2f1eeaded9b --position 1",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.IntVar(&s.rulePosition, "position", 0, "Rule position. Available: 1-1000")
	s.AddFlags(flagSet)

	commands.Must(s.Cobra().MarkFlagRequired("position"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("position", cobra.NoFileCompletions))
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if s.rulePosition < 1 || s.rulePosition > 1000 {
		return nil, fmt.Errorf("invalid position (1-1000 allowed)")
	}
	msg := fmt.Sprintf("Deleting firewall rule %d from server %v", s.rulePosition, arg)
	exec.PushProgressStarted(msg)

	err := exec.Firewall().DeleteFirewallRule(exec.Context(), &request.DeleteFirewallRuleRequest{
		ServerUUID: arg,
		Position:   s.rulePosition,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
