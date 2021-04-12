package networkinterface

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	bootable          string
	sourceIPfiltering string
	ipAddresses       []string
	newIndex          int
	currentIndex      int
	resolver.CachingServer
	completion.Server
}

// ModifyCommand creates the "network-interface modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify a network interface"),
	}
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	// TODO: reimplmement
	// s.SetPositionalArgHelp(server.PositionalArgHelp)
	fs := &pflag.FlagSet{}
	fs.IntVar(&s.currentIndex, "index", s.currentIndex, "Index of the interface to modify.")
	fs.IntVar(&s.newIndex, "new-index", s.newIndex, "Index of the interface to modify.")
	// TODO: refactor string to tristate bools (eg. allow empty)
	fs.StringVar(&s.bootable, "bootable", s.bootable, "Whether to try booting through the interface.")
	fs.StringVar(&s.sourceIPfiltering, "source-ip-filtering", s.sourceIPfiltering, "Whether source IP filtering is enabled on the interface. Disabling it is allowed only for SDN private interfaces.")
	fs.StringSliceVar(&s.ipAddresses, "ip-addresses", s.ipAddresses, "Array of IP addresses, multiple can be declared\nUsage: --ip-address address=94.237.112.143,family=IPv4")
	s.AddFlags(fs) // TODO(ana): replace usage with examples once the refactor is done.
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *modifyCommand) MaximumExecutions() int {
	return maxNetworkInterfaceActions
}

// Execute implements command.Command
func (s *modifyCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if s.currentIndex == 0 {
		return nil, fmt.Errorf("index is required")
	}
	ipAddresses, err := mapIPAddressesToRequest(s.ipAddresses)
	if err != nil {
		return nil, err
	}
	// initialize bootable and filtering flags as empty
	var empty = upcloud.Empty
	var bootable *upcloud.Boolean = &empty
	var sourceIPFiltering *upcloud.Boolean = &empty
	if s.bootable != "" {
		bootable, err = commands.BoolFromString(s.bootable)
		if err != nil {
			return nil, err
		}
	}
	if s.sourceIPfiltering != "" {
		sourceIPFiltering, err = commands.BoolFromString(s.sourceIPfiltering)
		if err != nil {
			return nil, err
		}
	}
	msg := fmt.Sprintf("Modifying network interface %q of server %q", s.currentIndex, arg)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	res, err := exec.Network().ModifyNetworkInterface(&request.ModifyNetworkInterfaceRequest{
		ServerUUID:        arg,
		CurrentIndex:      s.currentIndex,
		NewIndex:          s.newIndex,
		IPAddresses:       ipAddresses,
		SourceIPFiltering: *sourceIPFiltering,
		Bootable:          *bootable,
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
