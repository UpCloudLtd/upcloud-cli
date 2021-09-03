package networkinterface

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"

	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
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
		BaseCommand: commands.New(
			"modify",
			"Modify a network interface",
			"upctl server network-interface modify 009d7f4e-99ce-4c78-88f1-e695d4c37743 --index 2 --new-index 1",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.IntVar(&s.currentIndex, "index", s.currentIndex, "Index of the interface to modify.")
	fs.IntVar(&s.newIndex, "new-index", s.newIndex, "New index to move the interface to.")
	// TODO: refactor string to tristate bools (eg. allow empty)
	fs.StringVar(&s.bootable, "bootable", s.bootable, "Whether to try booting through the interface.")
	fs.StringVar(&s.sourceIPfiltering, "source-ip-filtering", s.sourceIPfiltering, "Whether source IP filtering is enabled on the interface. Disabling it is allowed only for SDN private interfaces.")
	fs.StringSliceVar(&s.ipAddresses, "ip-addresses", s.ipAddresses, "A comma-separated list of IP addresses, multiple can be declared\nUsage: --ip-address address=94.237.112.143,family=IPv4")
	s.AddFlags(fs) // TODO(ana): replace usage with examples once the refactor is done.
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *modifyCommand) MaximumExecutions() int {
	return maxNetworkInterfaceActions
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *modifyCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	if s.currentIndex == 0 {
		return nil, fmt.Errorf("index is required")
	}
	ipAddresses, err := mapIPAddressesToRequest(s.ipAddresses)
	if err != nil {
		return nil, err
	}
	// initialize bootable and filtering flags as empty upcloud.Booleans
	var empty = upcloud.Empty
	var bootable = &empty
	var sourceIPFiltering = &empty
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

	return output.OnlyMarshaled{Value: res}, nil
}
