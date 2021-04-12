package serverfirewall

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
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
		BaseCommand: commands.New("delete", "Removes a firewall rule from a server. Firewall rules must be removed individually. The positions of remaining firewall rules will be adjusted after a rule is removed."),
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.IntVar(&s.rulePosition, "position", 0, "Rule position. Available: 1-1000")
	s.AddFlags(flagSet)
}

// Execute implements Command.Execute
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if arg == "" {
		return nil, fmt.Errorf("server is required")
	}
	if s.rulePosition == 0 {
		return nil, fmt.Errorf("position is required")
	}
	if s.rulePosition < 1 || s.rulePosition > 1000 {
		return nil, fmt.Errorf("invalid position (1-1000 allowed)")
	}
	msg := fmt.Sprintf("deleting firewall rule %d from server %v", s.rulePosition, arg)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	err := exec.Firewall().DeleteFirewallRule(&request.DeleteFirewallRuleRequest{
		ServerUUID: arg,
		Position:   s.rulePosition,
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.None{}, nil
}
