package ipaddress

import (
	"testing"

	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
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

func TestSearchStorage(t *testing.T) {

	var IPAddress1 = upcloud.IPAddress{
		Address:   "94.237.113.140",
		PTRRecord: "ptr-record-1",
	}

	var IPAddress2 = upcloud.IPAddress{
		Address:   "94.237.113.141",
		PTRRecord: "ptr-record-2",
	}

	var IPAddress3 = upcloud.IPAddress{
		Address:   "94.237.113.142",
		PTRRecord: "ptr-record-3",
	}

	var IPAddress4 = upcloud.IPAddress{
		Address:   "94.237.113.143",
		PTRRecord: IPAddress1.PTRRecord,
	}

	var IPAddresses = upcloud.IPAddresses{IPAddresses: []upcloud.IPAddress{
		IPAddress1,
		IPAddress2,
		IPAddress3,
		IPAddress4,
	}}

	for _, testcase := range []struct {
		name         string
		args         []string
		expected     []string
		unique       bool
		additional   []upcloud.Storage
		backendCalls int
		errMsg       string
	}{
		{
			name:         "SingleAddress",
			args:         []string{IPAddress2.Address},
			expected:     []string{IPAddress2.Address},
			backendCalls: 0,
		},
		{
			name:         "MultipleAddressSearched",
			args:         []string{IPAddress2.Address, IPAddress3.Address},
			expected:     []string{IPAddress2.Address, IPAddress3.Address},
			backendCalls: 0,
		},
		{
			name:         "SinglePTRRecord",
			args:         []string{IPAddress2.PTRRecord},
			expected:     []string{IPAddress2.Address},
			backendCalls: 1,
		},
		{
			name:         "MultiplePTRRecordsSearched",
			args:         []string{IPAddress2.PTRRecord, IPAddress3.PTRRecord},
			expected:     []string{IPAddress2.Address, IPAddress3.Address},
			backendCalls: 1,
		},
		{
			name:         "MultiplePTRRecordsFound",
			args:         []string{IPAddress1.PTRRecord},
			expected:     []string{IPAddress1.Address, IPAddress4.Address},
			backendCalls: 1,
		},
		{
			name:         "MultiplePTRRecordsFound_UniqueWanted",
			args:         []string{IPAddress1.PTRRecord},
			expected:     []string{IPAddress1.Address, IPAddress4.Address},
			backendCalls: 1,
			unique:       true,
			errMsg:       "multiple ip addresses matched to query \"" + IPAddress1.PTRRecord + "\", use Address to specify",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			cachedIPs = nil
			mService := smock.Service{}
			mService.On("GetIPAddresses").Return(&IPAddresses, nil)

			result, err := searchIPAddresses(testcase.args, &mService, testcase.unique)

			if testcase.errMsg == "" {
				assert.Nil(t, err)
				assert.ElementsMatch(t, testcase.expected, result)
			} else {
				assert.Nil(t, result)
				assert.EqualError(t, err, testcase.errMsg)
			}
			mService.AssertNumberOfCalls(t, "GetIPAddresses", testcase.backendCalls)
		})
	}
}
