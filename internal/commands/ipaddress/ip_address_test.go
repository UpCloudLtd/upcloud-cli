package ipaddress

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestGetFamily(t *testing.T) {
	for _, test := range []struct {
		name     string
		address  string
		expected string
	}{
		{
			name:     "valid IPv4",
			address:  "127.0.0.1",
			expected: upcloud.IPAddressFamilyIPv4,
		},
		{
			name:     "valid IPv4 CIDR",
			address:  "127.0.0.1/24",
			expected: upcloud.IPAddressFamilyIPv4,
		},
		{
			name:     "invalid IPv4",
			address:  "127.0.0.300",
			expected: "127.0.0.300 is an invalid ip address",
		},
		{
			name:     "valid IPv6",
			address:  "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expected: upcloud.IPAddressFamilyIPv6,
		},
		{
			name:     "valid IPv6 CIDR",
			address:  "2001:0db8:85a3:0000:0000:8a2e:0370:7334/32",
			expected: upcloud.IPAddressFamilyIPv6,
		},
		{
			name:     "invalid IPv6",
			address:  "2001:0db8:85a3:0000:0000:8a2e:0370:g539",
			expected: "2001:0db8:85a3:0000:0000:8a2e:0370:g539 is an invalid ip address",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			family, err := GetFamily(test.address)
			if err != nil {
				assert.Equal(t, test.expected, err.Error())
			} else {
				assert.Equal(t, test.expected, family)
			}
		})
	}
}
