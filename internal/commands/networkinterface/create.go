package networkinterface

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	ipAddresses       []string
	bootable          config.OptionalBoolean
	sourceIPFiltering config.OptionalBoolean
	networkUUID       string
	family            string
	interfaceIndex    int
	networkType       string
	resolver.CachingServer
	completion.Server
}

// CreateCommand creates the "network-interface create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a network interface",
			"upctl server network-interface create 009d7f4e-99ce-4c78-88f1-e695d4c37743 --network 037a530b-533e-4cef-b6ad-6af8094bb2bc",
			"upctl server network-interface create 009d7f4e-99ce-4c78-88f1-e695d4c37743 --type public --source-ip-filtering",
			"upctl server network-interface create 009d7f4e-99ce-4c78-88f1-e695d4c37743 --type public --source-ip-filtering --ip-addresses 94.237.112.143,94.237.112.144",
			"upctl server network-interface create my_server2 --network 037a530b-533e-4cef-b6ad-6af8094bb2bc",
		),
	}
}

const (
	defaultNetworkType     = upcloud.NetworkTypePrivate
	defaultIPAddressFamily = upcloud.IPAddressFamilyIPv4
)

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.networkUUID, "network", "", "Virtual network ID or name to join.")
	fs.StringVar(&s.networkType, "type", defaultNetworkType, "Set the type of the network. Available: public, utility, private")
	fs.StringVar(&s.family, "family", defaultIPAddressFamily, "The address family of new IP address.")
	fs.IntVar(&s.interfaceIndex, "index", 0, "Interface index.")
	config.AddEnableOrDisableFlag(fs, &s.bootable, false, "bootable", "booting through the interface")
	config.AddEnableOrDisableFlag(fs, &s.sourceIPFiltering, false, "source-ip-filtering", "filtering source IPs")
	fs.StringSliceVar(&s.ipAddresses, "ip-addresses", []string{}, "A comma-separated list of IP addresses")
	s.AddFlags(fs)
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *createCommand) MaximumExecutions() int {
	return maxNetworkInterfaceActions
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *createCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	ipAddresses := []request.CreateNetworkInterfaceIPAddress{}
	if s.networkUUID == "" {
		ipAddresses = request.CreateNetworkInterfaceIPAddressSlice{{Family: s.family}}
	} else {
		if len(s.ipAddresses) == 0 {
			ipFamily := upcloud.IPAddressFamilyIPv4
			// Currently only IPv4 is supported in private networks
			if s.family != "IPv4" && s.networkType == "private" {
				return nil, fmt.Errorf("currently only IPv4 is supported in private networks")
			}
			if s.family != "" {
				ipFamily = s.family
			}
			ip := request.CreateNetworkInterfaceIPAddress{
				Family: ipFamily,
			}
			ipAddresses = append(ipAddresses, ip)
		} else {
			handled, err := mapIPAddressesToRequest(s.ipAddresses)
			if err != nil {
				return nil, err
			}
			ipAddresses = handled
		}
		res, err := exec.Network().GetNetworkDetails(&request.GetNetworkDetailsRequest{UUID: s.networkUUID})
		if err != nil {
			return nil, fmt.Errorf("invalid network requested: %w", err)
		}
		s.networkUUID = res.UUID
	}

	msg := fmt.Sprintf("Creating network interface for server %s network %s", arg, s.networkUUID)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	res, err := exec.Network().CreateNetworkInterface(&request.CreateNetworkInterfaceRequest{
		ServerUUID:        arg,
		Type:              s.networkType,
		NetworkUUID:       s.networkUUID,
		Index:             s.interfaceIndex,
		IPAddresses:       ipAddresses,
		SourceIPFiltering: s.sourceIPFiltering.AsUpcloudBoolean(),
		Bootable:          s.bootable.AsUpcloudBoolean(),
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
