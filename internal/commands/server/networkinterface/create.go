package networkinterface

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/ipaddress"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/network"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	ipAddresses       []string
	bootable          config.OptionalBoolean
	sourceIPFiltering config.OptionalBoolean
	networkArg        string
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
			"upctl server network-interface create 009d7f4e-99ce-4c78-88f1-e695d4c37743 --type private --network 037a530b-533e-4cef-b6ad-6af8094bb2bc --disable-source-ip-filtering --ip-addresses 10.0.0.1",
			"upctl server network-interface create my_server2 --type public --family IPv6",
			"upctl server network-interface create my_server2 --type public --family IPv4",
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
	fs.StringVar(&s.networkArg, "network", "", "Private network name or UUID to join.")
	fs.StringVar(&s.networkType, "type", defaultNetworkType, "Set the type of the network. Available: "+strings.Join(network.Types, ", "))
	fs.StringVar(&s.family, "family", defaultIPAddressFamily, "The address family of new IP address.")
	fs.IntVar(&s.interfaceIndex, "index", 0, "Interface index.")
	config.AddEnableDisableFlags(fs, &s.bootable, "bootable", "Whether to try booting through the interface.")
	config.AddEnableDisableFlags(fs, &s.sourceIPFiltering, "source-ip-filtering", "Whether source IP filtering is enabled on the interface. Disabling it is allowed only for SDN private interfaces.")
	fs.StringSliceVar(&s.ipAddresses, "ip-addresses", []string{}, "A comma-separated list of IP addresses")
	s.AddFlags(fs)
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("type", cobra.FixedCompletions(network.Types, cobra.ShellCompDirectiveNoFileComp)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("family", cobra.FixedCompletions(ipaddress.Families, cobra.ShellCompDirectiveNoFileComp)))
	for _, flag := range []string{"index", "ip-addresses"} {
		commands.Must(s.Cobra().RegisterFlagCompletionFunc(flag, cobra.NoFileCompletions))
	}
}

func (s *createCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("network", namedargs.CompletionFunc(completion.Network{}, cfg)))
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *createCommand) MaximumExecutions() int {
	return maxNetworkInterfaceActions
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *createCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	ipAddresses := []request.CreateNetworkInterfaceIPAddress{}
	if len(s.ipAddresses) == 0 {
		ipFamily := defaultIPAddressFamily
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

	// Resolve network UUID if --network arg was given
	networkUUID := ""
	if s.networkArg != "" {
		var err error
		networkUUID, err = namedargs.ResolveNetwork(exec, s.networkArg)
		if err != nil {
			return nil, err
		}
	}

	msg := fmt.Sprintf("Creating network interface for server %s to network %s", arg, networkUUID)
	exec.PushProgressStarted(msg)

	res, err := exec.Network().CreateNetworkInterface(exec.Context(), &request.CreateNetworkInterfaceRequest{
		ServerUUID:        arg,
		Type:              s.networkType,
		NetworkUUID:       networkUUID,
		Index:             s.interfaceIndex,
		IPAddresses:       ipAddresses,
		SourceIPFiltering: s.sourceIPFiltering.AsUpcloudBoolean(),
		Bootable:          s.bootable.AsUpcloudBoolean(),
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "MAC Address", Value: res.MAC},
		{Title: "IP Addresses", Value: res, Format: formatIPAddresses},
	}}, nil
}

func formatIPAddresses(val interface{}) (text.Colors, string, error) {
	iface, ok := val.(*upcloud.Interface)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse IP addresses from %T, expected *upcloud.Interface", val)
	}

	strs := make([]string, len(iface.IPAddresses))

	for i, ipa := range iface.IPAddresses {
		strs[i] = ui.DefaultAddressColours.Sprint(ipa.Address)
	}

	return nil, strings.Join(strs, ",\n"), nil
}
