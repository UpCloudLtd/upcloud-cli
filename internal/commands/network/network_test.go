package network

import (
	"testing"

	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchAllNetworks(t *testing.T) {

	var Network1 = upcloud.Network{
		Name: "network-1",
		UUID: "28e15cf5-8817-42ab-b017-970666be96ec",
	}

	var Network2 = upcloud.Network{
		Name: "network-2",
		UUID: "f9f5ad16-a63a-4670-8449-c01d1e97281e",
	}

	var Network3 = upcloud.Network{
		Name: "network-3",
		UUID: "e157ce0a-eeb0-49fc-9f2c-a05c3ac57066",
	}

	var Network4 = upcloud.Network{
		Name: Network1.Name,
		UUID: "b3e49768-f13a-42c3-bea7-4e2471657f2f",
	}

	var networks = &upcloud.Networks{Networks: []upcloud.Network{Network1, Network2, Network3, Network4}}

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
			name:         "SingleUuid",
			args:         []string{Network2.UUID},
			expected:     []string{Network2.UUID},
			backendCalls: 0,
		},
		{
			name:         "MultipleUuidSearched",
			args:         []string{Network2.UUID, Network3.UUID},
			expected:     []string{Network2.UUID, Network3.UUID},
			backendCalls: 0,
		},
		{
			name:         "SingleName",
			args:         []string{Network2.Name},
			expected:     []string{Network2.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleNamesSearched",
			args:         []string{Network2.Name, Network3.Name},
			expected:     []string{Network2.UUID, Network3.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleNamesFound",
			args:         []string{Network1.Name},
			expected:     []string{Network1.UUID, Network4.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleNamesFound_UniqueWanted",
			args:         []string{Network1.Name},
			expected:     []string{Network1.UUID, Network4.UUID},
			backendCalls: 1,
			unique:       true,
			errMsg:       "multiple networks matched to query \"" + Network1.Name + "\", use UUID to specify",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			cachedNetworks = nil
			mService := smock.MockService{}
			mService.On("GetNetworks", mock.Anything).Return(networks, nil)

			result, err := searchAllNetworks(testcase.args, &mService, testcase.unique)

			if testcase.errMsg == "" {
				assert.Nil(t, err)
				assert.ElementsMatch(t, testcase.expected, result)
			} else {
				assert.Nil(t, result)
				assert.EqualError(t, err, testcase.errMsg)
			}
			mService.AssertNumberOfCalls(t, "GetNetworks", testcase.backendCalls)
		})
	}
}

func TestSearchSUniqueNetwork(t *testing.T) {

	var Network1 = upcloud.Network{
		Name: "network-1",
		UUID: "28e15cf5-8817-42ab-b017-970666be96ec",
	}

	var Network2 = upcloud.Network{
		Name: "network-2",
		UUID: "f9f5ad16-a63a-4670-8449-c01d1e97281e",
	}

	var Network3 = upcloud.Network{
		Name: "network-3",
		UUID: "e157ce0a-eeb0-49fc-9f2c-a05c3ac57066",
	}

	var Network4 = upcloud.Network{
		Name: Network1.Name,
		UUID: "b3e49768-f13a-42c3-bea7-4e2471657f2f",
	}

	var networks = &upcloud.Networks{Networks: []upcloud.Network{Network1, Network2, Network3, Network4}}

	for _, testcase := range []struct {
		name         string
		args         string
		expected     *upcloud.Network
		backendCalls int
		errMsg       string
	}{
		{
			name:         "SingleUuid",
			args:         Network2.UUID,
			expected:     &Network2,
			backendCalls: 1,
		},
		{
			name:         "SingleName",
			args:         Network2.Name,
			expected:     &Network2,
			backendCalls: 1,
		},
		{
			name:         "MultipleNamesFound_UniqueWanted",
			args:         Network1.Name,
			backendCalls: 1,
			errMsg:       "multiple networks matched to query \"" + Network1.Name + "\", use UUID to specify",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			cachedNetworks = nil
			mService := smock.MockService{}
			mService.On("GetNetworks", mock.Anything).Return(networks, nil)

			result, err := SearchUniqueNetwork(testcase.args, &mService)

			if testcase.errMsg == "" {
				assert.Nil(t, err)
				assert.Equal(t, testcase.expected, result)
			} else {
				assert.Nil(t, result)
				assert.EqualError(t, err, testcase.errMsg)
			}
			mService.AssertNumberOfCalls(t, "GetNetworks", testcase.backendCalls)
		})
	}
}
