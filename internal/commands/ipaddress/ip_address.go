package ipaddress

import (
	"fmt"
	"net"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
)

const maxIPAddressActions = 10

// BaseIPAddressCommand creates the base 'ip-address' command
func BaseIPAddressCommand() commands.Command {
	return &ipAddressCommand{
		commands.New("ip-address", "Manage ip address"),
	}
}

type ipAddressCommand struct {
	*commands.BaseCommand
}

// GetFamily returns a human-readable string (`"IPv4"` or `"IPv6"`) of the address family of the ip parsed from the string
func GetFamily(address string) (string, error) {
	parsed := net.ParseIP(address)
	if parsed.To4() != nil {
		return upcloud.IPAddressFamilyIPv4, nil
	}
	if parsed.To16() != nil {
		return upcloud.IPAddressFamilyIPv6, nil
	}
	ip, _, err := net.ParseCIDR(address)
	if err != nil {
		return "", fmt.Errorf("%s is an invalid ip address", address)
	}
	if ip.To4() != nil {
		return upcloud.IPAddressFamilyIPv4, nil
	}
	if ip.To16() != nil {
		return upcloud.IPAddressFamilyIPv6, nil
	}
	return "", fmt.Errorf("%s is an invalid ip address", address)
}
