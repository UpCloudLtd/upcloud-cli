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

var mockServers = &upcloud.Servers{Servers: []upcloud.Server{
	{Title: "mock1", UUID: "abcdef", Hostname: "foo"},
	{Title: "mock2", UUID: "abcghi", Hostname: "afoo"},
	{Title: "bock1", UUID: "jklmno", Hostname: "faa"},
	{Title: "bock2", UUID: "pqrstu", Hostname: "fii"},
	{Title: "dock1", UUID: "vwxyzä", Hostname: "bfoo"},
	{Title: "case started", UUID: "started", Hostname: "astarted", State: "started"},
	{Title: "case stopped", UUID: "stopped", Hostname: "astopped", State: "stopped"},
}}

func TestServer_CompleteArgument(t *testing.T) {
	mService := new(smock.Service)
	mService.On("GetServers", mock.Anything).Return(mockServers, nil)

	for _, test := range []struct {
		name              string
		complete          string
		expectedMatches   []string
		expectedDirective cobra.ShellCompDirective
	}{
		{name: "basic uuid", complete: "pqr", expectedMatches: []string{"pqrstu"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "basic title", complete: "dock", expectedMatches: []string{"dock1"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "basic hostname", complete: "fa", expectedMatches: []string{"faa"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple uuids", complete: "abc", expectedMatches: []string{"abcdef", "abcghi"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple titles", complete: "bock", expectedMatches: []string{"bock1", "bock2"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple hostnames", complete: "f", expectedMatches: []string{"foo", "faa", "fii"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "hostnames and titles", complete: "b", expectedMatches: []string{"bock1", "bock2", "bfoo"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
	} {
		t.Run(test.name, func(t *testing.T) {
			servers, directive := completion.Server{}.CompleteArgument(context.TODO(), mService, test.complete)
			assert.Equal(t, test.expectedMatches, servers)
			assert.Equal(t, test.expectedDirective, directive)
		})
	}

	t.Run("stopped", func(t *testing.T) {
		servers, directive := completion.StoppedServer{}.CompleteArgument(context.TODO(), mService, "s")
		assert.Equal(t, []string{"stopped"}, servers)
		assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
	})

	t.Run("started", func(t *testing.T) {
		servers, directive := completion.StartedServer{}.CompleteArgument(context.TODO(), mService, "s")
		assert.Equal(t, []string{"started"}, servers)
		assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
	})
}

func TestServer_CompleteArgumentServiceFail(t *testing.T) {
	mService := new(smock.Service)
	mService.On("GetServers", mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))
	ips, directive := completion.Server{}.CompleteArgument(context.TODO(), mService, "127")
	assert.Nil(t, ips)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}
