package ipaddress

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type assignCommand struct {
	*commands.BaseCommand
	ipSvc     service.IpAddress
	serverSvc service.Server
	req       request.AssignIPAddressRequest
	floating  bool
}

// AssignCommand creates the 'ip-address assign' command
func AssignCommand(serverSvc service.Server, ipSvc service.IpAddress) commands.Command {
	return &assignCommand{
		BaseCommand: commands.New("assign", "Assign an ip address"),
		serverSvc:   serverSvc,
		ipSvc:       ipSvc,
	}
}

var defCreateParams = request.AssignIPAddressRequest{
	Access: upcloud.IPAddressAccessPublic,
	Family: upcloud.IPAddressFamilyIPv4,
}

// InitCommand implements Command.InitCommand
func (s *assignCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.req.Access, "access", defCreateParams.Access, "Is address for utility or public network.")
	fs.StringVar(&s.req.Family, "family", defCreateParams.Family, "The address family of new IP address.")
	fs.StringVar(&s.req.ServerUUID, "server", defCreateParams.ServerUUID, "The server the ip address is assigned to.")
	fs.StringVar(&s.req.MAC, "mac", defCreateParams.MAC, "MAC address of server interface to assign address to. Required for non-floating addresses.")
	fs.StringVar(&s.req.Zone, "zone", defCreateParams.Zone, "Zone of address, required when assigning a detached floating IP address.")
	fs.BoolVar(&s.floating, "floating", false, "Whether the address to be assigned is a floating one.")
	s.AddFlags(fs)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *assignCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		if s.floating && s.req.Zone == "" && s.req.MAC == "" {
			return nil, fmt.Errorf("MAC or zone is required for floating IP")
		}
		if !s.floating && s.req.ServerUUID == "" {
			return nil, fmt.Errorf("server is required for non-floating IP")
		}
		s.req.Floating = upcloud.FromBool(s.floating)

		if s.req.ServerUUID != "" {
			svr, err := server.SearchSingleServer(s.req.ServerUUID, s.serverSvc)
			if err != nil {
				return nil, err
			}
			s.req.ServerUUID = svr.UUID
		}

		return ui.HandleContext{
			RequestID: func(in interface{}) string {
				req := in.(*request.AssignIPAddressRequest)
				if req.MAC != "" {
					return req.MAC
				}
				return req.Zone
			},
			ResultUUID:    func(in interface{}) string { return in.(*upcloud.IPAddress).Address },
			ResultPrefix:  "IP Address",
			MaxActions:    maxIPAddressActions,
			InteractiveUI: s.Config().InteractiveUI(),
			ActionMsg:     "Assigning IP Address to",
			Action: func(req interface{}) (interface{}, error) {
				return s.ipSvc.AssignIPAddress(req.(*request.AssignIPAddressRequest))
			},
		}.Handle(commands.ToArray(&s.req))
	}
}