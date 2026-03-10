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

func TestObjectStorage_CompleteArgument(t *testing.T) {
	mService := new(smock.Service)
	mService.On("GetManagedObjectStorages", mock.Anything, mock.Anything).Return([]upcloud.ManagedObjectStorage{
		{UUID: "objsto-uuid-1", Name: "MockBucket"},
		{UUID: "objsto-uuid-2", Name: "another-bucket"},
	}, nil)

	vals, directive := completion.ObjectStorage{}.CompleteArgument(context.TODO(), mService, "MOCK")
	assert.Equal(t, []string{"MockBucket"}, vals)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}

func TestObjectStorage_CompleteArgumentServiceFail(t *testing.T) {
	mService := new(smock.Service)
	mService.On("GetManagedObjectStorages", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("MOCKFAIL"))

	vals, directive := completion.ObjectStorage{}.CompleteArgument(context.TODO(), mService, "mock")
	assert.Nil(t, vals)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}
