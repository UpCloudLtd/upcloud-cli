package ipaddress

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/spf13/pflag"
)

type assignCommand struct {
	*commands.BaseCommand
	floating   config.OptionalBoolean
	access     string
	family     string
	serverUUID string
	mac        string
	zone       string
}

// AssignCommand creates the 'ip-address assign' command
func AssignCommand() commands.Command {
	return &assignCommand{
		BaseCommand: commands.New(
			"assign",
			"Assign an ip address",
			"upctl ip-address assign --server 00038afc-d526-4148-af0e-d2f1eeaded9b",
			"upctl ip-address assign --server 00944977-89ce-4d10-89c3-bb5ba482e48d --family IPv6",
			"upctl ip-address assign --server 00944977-89ce-4d10-89c3-bb5ba482e48d --floating --zone pl-waw1",
			"upctl ip-address assign --server 00b78f8b-521d-4ffb-8baa-adf96c7b8f45 --floating --mac d6:0e:4a:6f:11:8f",
		),
	}
}

const (
	defaultAccess = upcloud.IPAddressAccessPublic
	defaultFamily = upcloud.IPAddressFamilyIPv4
)

// InitCommand implements Command.InitCommand
func (s *assignCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.access, "access", defaultAccess, "Is address for utility or public network.")
	fs.StringVar(&s.family, "family", defaultFamily, "The address family of new IP address.")
	fs.StringVar(&s.serverUUID, "server", "", "The server the ip address is assigned to.")
	fs.StringVar(&s.mac, "mac", "", "MAC address of server interface to assign address to. Required for detached floating IP address if zone is not specified.")
	fs.StringVar(&s.zone, "zone", "", "Zone of address. Required for detached floating IP address if MAC address is not speficied.")
	config.AddToggleFlag(fs, &s.floating, "floating", false, "Whether the address to be assigned is a floating one.")
	s.AddFlags(fs)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *assignCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	if s.floating.Value() && s.zone == "" && s.mac == "" {
		return nil, fmt.Errorf("MAC or zone is required for floating IP")
	}
	if !s.floating.Value() && s.serverUUID == "" {
		return nil, fmt.Errorf("server is required for non-floating IP")
	}

	if s.serverUUID != "" {
		_, err := exec.Server().GetServerDetails(&request.GetServerDetailsRequest{UUID: s.serverUUID})
		if err != nil {
			return nil, fmt.Errorf("invalid server uuid: %w", err)
		}
	}
	target := s.mac
	if target == "" {
		target = s.zone
	}
	msg := fmt.Sprintf("Assigning IP Address to %v", target)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	res, err := exec.IPAddress().AssignIPAddress(&request.AssignIPAddressRequest{
		Access:     s.access,
		Family:     s.family,
		ServerUUID: s.serverUUID,
		Floating:   s.floating.AsUpcloudBoolean(),
		MAC:        s.mac,
		Zone:       s.zone,
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}
	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()
	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "IP Address", Value: res.Address, Colour: ui.DefaultAddressColours},
	}}, nil
}
