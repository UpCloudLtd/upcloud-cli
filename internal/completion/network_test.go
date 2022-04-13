package completion_test

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockNetworks = &upcloud.Networks{Networks: []upcloud.Network{
	{Name: "mock1", UUID: "abcdef"},
	{Name: "mock2", UUID: "abcghi"},
	{Name: "bock1", UUID: "jklmno"},
	{Name: "bock2", UUID: "pqrstu"},
	{Name: "dock1", UUID: "vwxyzä"},
}}

func TestNetwork_CompleteArgument(t *testing.T) {
	for _, test := range []struct {
		name              string
		complete          string
		expectedMatches   []string
		expectedDirective cobra.ShellCompDirective
	}{
		{name: "basic uuid", complete: "pqr", expectedMatches: []string{"pqrstu"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "basic name", complete: "dock", expectedMatches: []string{"dock1"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple uuids", complete: "abc", expectedMatches: []string{"abcdef", "abcghi"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple names", complete: "bock", expectedMatches: []string{"bock1", "bock2"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := new(smock.Service)
			mService.On("GetNetworks", mock.Anything).Return(mockNetworks, nil)
			ips, directive := completion.Network{}.CompleteArgument(mService, test.complete)
			assert.Equal(t, test.expectedMatches, ips)
			assert.Equal(t, test.expectedDirective, directive)
		})
	}
}

func TestNetwork_CompleteArgumentServiceFail(t *testing.T) {
	mService := new(smock.Service)
	mService.On("GetNetworks", mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))
	ips, directive := completion.Network{}.CompleteArgument(mService, "127")
	assert.Nil(t, ips)
	assert.Equal(t, cobra.ShellCompDirectiveDefault, directive)
}
