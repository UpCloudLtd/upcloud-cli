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

var mockStorages = &upcloud.Storages{Storages: []upcloud.Storage{
	{Title: "mock1", UUID: "abcdef"},
	{Title: "mock2", UUID: "abcghi"},
	{Title: "bock1", UUID: "jklmno"},
	{Title: "bock2", UUID: "pqrstu"},
	{Title: "dock1", UUID: "vwxyzä"},
}}

func TestStorage_CompleteArgument(t *testing.T) {
	for _, test := range []struct {
		name              string
		complete          string
		expectedMatches   []string
		expectedDirective cobra.ShellCompDirective
	}{
		{name: "basic uuid", complete: "pqr", expectedMatches: []string{"pqrstu"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "basic title", complete: "dock", expectedMatches: []string{"dock1"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple uuids", complete: "abc", expectedMatches: []string{"abcdef", "abcghi"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple titles", complete: "bock", expectedMatches: []string{"bock1", "bock2"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := new(smock.Service)
			mService.On("GetStorages", mock.Anything).Return(mockStorages, nil)
			ips, directive := completion.Storage{}.CompleteArgument(mService, test.complete)
			assert.Equal(t, test.expectedMatches, ips)
			assert.Equal(t, test.expectedDirective, directive)
		})
	}
}

func TestStorage_CompleteArgumentServiceFail(t *testing.T) {
	mService := new(smock.Service)
	mService.On("GetStorages", mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))
	ips, directive := completion.Storage{}.CompleteArgument(mService, "127")
	assert.Nil(t, ips)
	assert.Equal(t, cobra.ShellCompDirectiveDefault, directive)
}
