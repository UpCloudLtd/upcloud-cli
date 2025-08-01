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

var mockClusters = []upcloud.KubernetesCluster{
	{Name: "asd-1", UUID: "abcdef"},
	{Name: "asd-2", UUID: "abcghi"},
	{Name: "qwe-1", UUID: "jklmno"},
}

func TestKubernetes_CompleteArgument(t *testing.T) {
	for _, test := range []struct {
		name              string
		complete          string
		expectedMatches   []string
		expectedDirective cobra.ShellCompDirective
	}{
		{name: "Name/UUID - no match", complete: "pqr", expectedMatches: []string(nil), expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "UUID - single match", complete: "jkl", expectedMatches: []string{"jklmno"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "UUID - multiple matches", complete: "abc", expectedMatches: []string{"abcdef", "abcghi"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "Name - one match", complete: "qwe", expectedMatches: []string{"qwe-1"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
		{name: "Name - multiple matches", complete: "asd", expectedMatches: []string{"asd-1", "asd-2"}, expectedDirective: cobra.ShellCompDirectiveNoFileComp},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := new(smock.Service)
			mService.On("GetKubernetesClusters", mock.Anything).Return(mockClusters, nil)
			completions, directive := completion.Kubernetes{}.CompleteArgument(context.TODO(), mService, test.complete)
			assert.Equal(t, test.expectedMatches, completions)
			assert.Equal(t, test.expectedDirective, directive)
		})
	}
}

func TestKubernetes_CompleteArgumentServiceFail(t *testing.T) {
	mService := new(smock.Service)
	mService.On("GetKubernetesClusters", mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))
	completions, directive := completion.Kubernetes{}.CompleteArgument(context.TODO(), mService, "asd")
	assert.Nil(t, completions)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}
