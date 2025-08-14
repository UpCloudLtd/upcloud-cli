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

var mockDatabases = []upcloud.ManagedDatabase{
	{Title: "asd-1", UUID: "abcdef"},
	{Title: "asd-2", UUID: "abcghi"},
	{Title: "qwe-1", UUID: "jklmno"},
}

func TestDatabase_CompleteArgument(t *testing.T) {
	for _, test := range []struct {
		name              string
		complete          string
		expectedMatches   []string
		expectedDirective cobra.ShellCompDirective
	}{
		{name: "Title/UUID - no match", complete: "pqr", expectedMatches: []string(nil), expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "UUID - single match", complete: "jkl", expectedMatches: []string{"jklmno"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "UUID - multiple matches", complete: "abc", expectedMatches: []string{"abcdef", "abcghi"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "Title - one match", complete: "qwe", expectedMatches: []string{"qwe-1"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "Title - multiple matches", complete: "asd", expectedMatches: []string{"asd-1", "asd-2"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := new(smock.Service)
			mService.On("GetManagedDatabases", mock.Anything).Return(mockDatabases, nil)
			completions, directive := completion.Database{}.CompleteArgument(context.TODO(), mService, test.complete)
			assert.Equal(t, test.expectedMatches, completions)
			assert.Equal(t, test.expectedDirective, directive)
		})
	}
}

func TestDatabase_CompleteArgumentServiceFail(t *testing.T) {
	mService := new(smock.Service)
	mService.On("GetManagedDatabases", mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))
	completions, directive := completion.Database{}.CompleteArgument(context.TODO(), mService, "asd")
	assert.Nil(t, completions)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}
