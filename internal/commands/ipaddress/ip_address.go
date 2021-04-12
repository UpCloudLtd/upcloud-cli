package ipaddress

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"net"
)

const maxIPAddressActions = 10

// TODO: re-add
// const positionalArgHelp = "<Address/PTRRecord...>"

// BaseIPAddressCommand creates the base 'ip-address' command
func BaseIPAddressCommand() commands.Command {
	return &ipAddressCommand{commands.New("ip-address", "Manage ip address")}
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

var cachedIPs []upcloud.IPAddress

func searchIPAddress(term string, service service.IpAddress, unique bool) ([]*upcloud.IPAddress, error) {
	var result []*upcloud.IPAddress

	if len(cachedIPs) == 0 {
		ips, err := service.GetIPAddresses()
		if err != nil {
			return nil, err
		}
		cachedIPs = ips.IPAddresses
	}

	for _, i := range cachedIPs {
		ip := i
		if ip.Address == term || ip.PTRRecord == term {
			result = append(result, &ip)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no ip address was found with %s", term)
	}
	if len(result) > 1 && unique {
		return nil, fmt.Errorf("multiple ip addresses matched to query %q, use Address to specify", term)
	}
	return result, nil
}

func searchIPAddresses(terms []string, service service.IpAddress, unique bool) ([]string, error) {
	var result []string
	for _, term := range terms {
		_, err := GetFamily(term)
		if err == nil {
			result = append(result, term)
		} else {
			ip, err := searchIPAddress(term, service, unique)
			if err != nil {
				return nil, err
			}
			for _, i := range ip {
				result = append(result, i.Address)
			}
		}
	}
	return result, nil
}
