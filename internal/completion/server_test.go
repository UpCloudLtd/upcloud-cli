package completion_test

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
)

var mockServers = &upcloud.Servers{Servers: []upcloud.Server{
	{Title: "mock1", UUID: "abcdef", Hostname: "foo"},
	{Title: "mock2", UUID: "abcghi", Hostname: "afoo"},
	{Title: "bock1", UUID: "jklmno", Hostname: "faa"},
	{Title: "bock2", UUID: "pqrstu", Hostname: "fii"},
	{Title: "dock1", UUID: "vwxyzä", Hostname: "bfoo"},
}}

func TestServer_CompleteArgument(t *testing.T) {
	t.Parallel()
	for _, test := range []completionTest{
		{name: "basic uuid", complete: "pqr", expectedMatches: []string{"pqrstu"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "basic title", complete: "dock", expectedMatches: []string{"dock1"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "basic hostname", complete: "fa", expectedMatches: []string{"faa"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple uuids", complete: "abc", expectedMatches: []string{"abcdef", "abcghi"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple titles", complete: "bock", expectedMatches: []string{"bock1", "bock2"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple hostnames", complete: "f", expectedMatches: []string{"foo", "faa", "fii"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "hostnames and titles", complete: "b", expectedMatches: []string{"bock1", "bock2", "bfoo"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
	} {
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			testCompletion(t, "GetServers", mockServers, completion.Server{}, test.complete, test.expectedMatches, test.expectedDirective)
		})
	}
}

func TestServer_CompleteArgumentServiceFail(t *testing.T) {
	t.Parallel()
	mService := new(smock.Service)
	mService.On("GetServers", mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))
	ips, directive := completion.Server{}.CompleteArgument(mService, "127")
	assert.Nil(t, ips)
	assert.Equal(t, cobra.ShellCompDirectiveDefault, directive)
}
