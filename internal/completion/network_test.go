package completion_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockNetworks = &upcloud.Networks{Networks: []upcloud.Network{
	{Name: "mock1", UUID: "abcdef", Type: "private"},
	{Name: "mock2", UUID: "abcghi", Type: "private"},
	{Name: "mock3", UUID: "123qwe", Type: "public"},
	{Name: "mock4", UUID: "123zxc", Type: "utility"},
	{Name: "bock1", UUID: "jklmno", Type: "private"},
	{Name: "bock2", UUID: "pqrstu", Type: "private"},
	{Name: "dock1", UUID: "vwxyzä", Type: "private"},
}}

func TestNetwork_CompleteArgument(t *testing.T) {
	for _, test := range []struct {
		name              string
		complete          string
		expectedMatches   []string
		expectedDirective cobra.ShellCompDirective
	}{
		{name: "basic uuid", complete: "pqr", expectedMatches: []string{"pqrstu"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "public net", complete: "123q", expectedMatches: nil, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "utility net", complete: "123z", expectedMatches: nil, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "basic name", complete: "dock", expectedMatches: []string{"dock1"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple uuids", complete: "abc", expectedMatches: []string{"abcdef", "abcghi"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple names", complete: "bock", expectedMatches: []string{"bock1", "bock2"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := new(smock.Service)
			mService.On("GetNetworks", mock.Anything).Return(mockNetworks, nil)
			ips, directive := completion.Network{}.CompleteArgument(context.TODO(), mService, test.complete)
			assert.Equal(t, test.expectedMatches, ips)
			assert.Equal(t, test.expectedDirective, directive)
		})
	}
}

func TestNetwork_CompleteArgumentServiceFail(t *testing.T) {
	mService := new(smock.Service)
	mService.On("GetNetworks", mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))
	ips, directive := completion.Network{}.CompleteArgument(context.TODO(), mService, "127")
	assert.Nil(t, ips)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}
