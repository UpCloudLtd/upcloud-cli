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

var mockIPs = &upcloud.IPAddresses{IPAddresses: []upcloud.IPAddress{
	{Address: "127.0.0.1", PTRRecord: "localhost"},
	{Address: "127.0.0.2", PTRRecord: "localmost"},
	{Address: "128.0.0.3", PTRRecord: "focalhost"},
}}

func TestIPAddress_CompleteArgument(t *testing.T) {
	for _, test := range []struct {
		name              string
		complete          string
		expectedMatches   []string
		expectedDirective cobra.ShellCompDirective
	}{
		{name: "basic ptr", complete: "localh", expectedMatches: []string{"localhost"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "basic ip", complete: "128", expectedMatches: []string{"128.0.0.3"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple ptrs", complete: "local", expectedMatches: []string{"localhost", "localmost"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple ips", complete: "127", expectedMatches: []string{"127.0.0.1", "127.0.0.2"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := new(smock.Service)
			mService.On("GetIPAddresses", mock.Anything).Return(mockIPs, nil)
			ips, directive := completion.IPAddress{}.CompleteArgument(context.TODO(), mService, test.complete)
			assert.Equal(t, test.expectedMatches, ips)
			assert.Equal(t, test.expectedDirective, directive)
		})
	}
}

func TestIPAddress_CompleteArgumentServiceFail(t *testing.T) {
	mService := new(smock.Service)
	mService.On("GetIPAddresses", mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))
	ips, directive := completion.IPAddress{}.CompleteArgument(context.TODO(), mService, "127")
	assert.Nil(t, ips)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}
