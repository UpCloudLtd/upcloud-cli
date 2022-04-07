package network

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
)

type modifyCommand struct {
	*commands.BaseCommand
	networks     []string
	attachRouter string
	name         string
	detachRouter config.OptionalBoolean
	completion.Network
	resolver.CachingNetwork
	// routerResolver is used to support resolving names of routers to uuids
	routerResolver resolver.CachingRouter
}

// ModifyCommand creates the "network modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New(
			"modify",
			"Modify a network",
			"upctl network modify 037a530b-533e-4cef-b6ad-6af8094bb2bc --ip-network dhcp=false,family=IPv4",
			`upctl network modify "My Network" --name "My Super Network"`,
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.name, "name", "", "Set name of the private network.")
	fs.StringVar(&s.attachRouter, "router", "", "Attach a router to this network, specified by router UUID or router name.")
	config.AddToggleFlag(fs, &s.detachRouter, "detach-router", false, "Detach a router from this network.")
	fs.StringArrayVar(&s.networks, "ip-network", s.networks, "The ip network with modified values. \n\n"+
		"Fields \n"+
		"  family: string \n"+
		"  gateway: string \n"+
		"  dhcp: true/false \n"+
		"  dhcp-default-route: true/false \n"+
		"  dhcp-dns: array of strings")
	s.AddFlags(fs)
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *modifyCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	if s.attachRouter != "" && s.detachRouter == config.True {
		return nil, fmt.Errorf("ambiguous command, cannot detach and attach a router at the same time")
	}
	var networks []upcloud.IPNetwork
	for _, networkStr := range s.networks {
		network, err := handleNetwork(networkStr)
		if err != nil {
			return nil, err
		}
		if network.Family == "" {
			return nil, fmt.Errorf("family is required")
		}
		network.Address = ""
		networks = append(networks, *network)
	}

	msg := fmt.Sprintf("modifying network %v", arg)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()

	var network *upcloud.Network
	if s.name != "" || len(networks) > 0 {
		// we want to update name and/or networks
		logline.SetMessage(fmt.Sprintf("%s: sending modify request", msg))
		res, err := exec.Network().ModifyNetwork(&request.ModifyNetworkRequest{
			UUID:       arg,
			Name:       s.name,
			Zone:       "", // TODO: should this be implemented?
			IPNetworks: networks,
		})
		if err != nil {
			logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
			logline.SetDetails(err.Error(), "error: ")
			return nil, err
		}
		// store the result in order to return it
		network = res
	}

	if s.attachRouter != "" {
		routerResolver, err := s.routerResolver.Get(exec.All())
		if err != nil {
			return nil, fmt.Errorf("cannot get router resolver: %w", err)
		}
		routerUUID, err := routerResolver(s.attachRouter)
		if err != nil {
			return nil, fmt.Errorf("cannot get resolve router '%s': %w", s.attachRouter, err)
		}
		logline.SetMessage(fmt.Sprintf("%s: attaching router %s", msg, routerUUID))
		logline.SetDetails(routerUUID, "router UUID: ")
		err = exec.Network().AttachNetworkRouter(&request.AttachNetworkRouterRequest{
			NetworkUUID: arg,
			RouterUUID:  routerUUID,
		})
		if err != nil {
			logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
			logline.SetDetails(err.Error(), "error: ")
			return nil, fmt.Errorf("cannot attach router '%s': %w", s.attachRouter, err)
		}
		// update the stored result (if we have one) manually to avoid refetching later
		if network != nil {
			network.Router = routerUUID
		}
	} else if s.detachRouter == config.True {
		logline.SetMessage(fmt.Sprintf("%s: detaching router", msg))
		err := exec.Network().DetachNetworkRouter(&request.DetachNetworkRouterRequest{
			NetworkUUID: arg,
		})
		if err != nil {
			logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
			logline.SetDetails(err.Error(), "error: ")
			return nil, fmt.Errorf("cannot detach router '%s': %w", s.attachRouter, err)
		}
		// update the stored result (if we have one) manually to avoid refetching later
		if network != nil {
			network.Router = ""
		}
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()
	if network == nil {
		// if we're just detaching/attaching, we won't have network returned from the calls so re-fetch
		details, err := exec.Network().GetNetworkDetails(&request.GetNetworkDetailsRequest{UUID: arg})
		if err != nil {
			return nil, fmt.Errorf("cannot get network state: %w", err)
		}
		network = details
	}
	return output.OnlyMarshaled{Value: network}, nil
}
