package completion_test

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockIPs = &upcloud.IPAddresses{IPAddresses: []upcloud.IPAddress{
	{Address: "127.0.0.1", PTRRecord: "localhost"},
	{Address: "127.0.0.2", PTRRecord: "localmost"},
	{Address: "128.0.0.3", PTRRecord: "focalhost"},
}}

func TestIPAddress_CompleteArgument(t *testing.T) {
	t.Parallel()
	for _, test := range []completionTest{
		{name: "basic ptr", complete: "localh", expectedMatches: []string{"localhost"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "basic ip", complete: "128", expectedMatches: []string{"128.0.0.3"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple ptrs", complete: "local", expectedMatches: []string{"localhost", "localmost"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple ips", complete: "127", expectedMatches: []string{"127.0.0.1", "127.0.0.2"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
	} {
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			testCompletion(t, "GetIPAddresses", mockIPs, completion.IPAddress{}, test.complete, test.expectedMatches, test.expectedDirective)
		})
	}
}

func TestIPAddress_CompleteArgumentServiceFail(t *testing.T) {
	t.Parallel()
	mService := new(smock.Service)
	mService.On("GetIPAddresses", mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))
	ips, directive := completion.IPAddress{}.CompleteArgument(mService, "127")
	assert.Nil(t, ips)
	assert.Equal(t, cobra.ShellCompDirectiveDefault, directive)
}
