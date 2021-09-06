package completion_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
)

type completionTest struct {
	name              string
	complete          string
	expectedMatches   []string
	expectedDirective cobra.ShellCompDirective
}

func testCompletion(t *testing.T, methodName string, returnValue interface{}, provider completion.Provider, toComplete string, expectedMatches []string, expectedDirective cobra.ShellCompDirective) {
	t.Helper()
	mService := new(smock.Service)
	mService.On(methodName, mock.Anything).Return(returnValue, nil)
	ips, directive := provider.CompleteArgument(mService, toComplete)
	assert.Equal(t, expectedMatches, ips)
	assert.Equal(t, expectedDirective, directive)
}
