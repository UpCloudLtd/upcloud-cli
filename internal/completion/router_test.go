//nolint:dupl // completion tests look *very* similar, but this is a false positive
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

var mockRouters = &upcloud.Routers{Routers: []upcloud.Router{
	{Name: "mock1", UUID: "abcdef"},
	{Name: "mock2", UUID: "abcghi"},
	{Name: "bock1", UUID: "jklmno"},
	{Name: "bock2", UUID: "pqrstu"},
	{Name: "dock1", UUID: "vwxyzä"},
}}

func TestRouter_CompleteArgument(t *testing.T) {
	t.Parallel()
	for _, test := range []completionTest{
		{name: "basic uuid", complete: "pqr", expectedMatches: []string{"pqrstu"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "basic name", complete: "dock", expectedMatches: []string{"dock1"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple uuids", complete: "abc", expectedMatches: []string{"abcdef", "abcghi"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple names", complete: "bock", expectedMatches: []string{"bock1", "bock2"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
	} {
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			testCompletion(t, "GetRouters", mockRouters, completion.Router{}, test.complete, test.expectedMatches, test.expectedDirective)
		})
	}
}

func TestRouter_CompleteArgumentServiceFail(t *testing.T) {
	t.Parallel()
	mService := new(smock.Service)
	mService.On("GetRouters", mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))
	ips, directive := completion.Router{}.CompleteArgument(mService, "127")
	assert.Nil(t, ips)
	assert.Equal(t, cobra.ShellCompDirectiveDefault, directive)
}
