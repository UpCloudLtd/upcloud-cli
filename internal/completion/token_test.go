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

var mockTokens = &upcloud.Tokens{
	{Name: "mock1", ID: "0c1eadbe-efde-adbe-efde-adbeefdeadbe"},
	{Name: "mock2", ID: "0c2eadbe-efde-adbe-efde-adbeefdeadbe"},
	{Name: "bock1", ID: "0c3eadbe-efde-adbe-efde-adbeefdeadbe"},
	{Name: "bock2", ID: "0c4eadbe-efde-adbe-efde-adbeefdeadbe"},
	{Name: "dock1", ID: "0c5eadbe-efde-adbe-efde-adbeefdeadbe"},
}

func TestToken_CompleteArgument(t *testing.T) {
	for _, test := range []struct {
		name              string
		complete          string
		expectedMatches   []string
		expectedDirective cobra.ShellCompDirective
	}{
		{name: "basic id", complete: "0c2", expectedMatches: []string{"0c2eadbe-efde-adbe-efde-adbeefdeadbe"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "multiple ids", complete: "0c", expectedMatches: []string{
			"0c1eadbe-efde-adbe-efde-adbeefdeadbe",
			"0c2eadbe-efde-adbe-efde-adbeefdeadbe",
			"0c3eadbe-efde-adbe-efde-adbeefdeadbe",
			"0c4eadbe-efde-adbe-efde-adbeefdeadbe",
			"0c5eadbe-efde-adbe-efde-adbeefdeadbe",
		}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := new(smock.Service)
			mService.On("GetTokens", mock.Anything).Return(mockTokens, nil)
			tokens, directive := completion.Token{}.CompleteArgument(context.TODO(), mService, test.complete)
			assert.Equal(t, test.expectedMatches, tokens)
			assert.Equal(t, test.expectedDirective, directive)
		})
	}
}

func TestToken_CompleteArgumentServiceFail(t *testing.T) {
	mService := new(smock.Service)
	mService.On("GetTokens", mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))
	tokens, directive := completion.Token{}.CompleteArgument(context.TODO(), mService, "FOO")
	assert.Nil(t, tokens)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}
