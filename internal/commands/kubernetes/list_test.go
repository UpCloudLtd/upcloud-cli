package kubernetes

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/testutils"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestList_MarshaledOutput(t *testing.T) {
	mService := smock.Service{}
	mService.On("GetKubernetesClusters", mock.Anything).Return([]upcloud.KubernetesCluster{testCluster}, nil)

	conf := config.New()
	conf.Viper().Set(config.KeyOutput, config.ValueOutputJSON)

	command := commands.BuildCommand(ListCommand(), nil, conf)

	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.Nil(t, err)
	testutils.AssertOutputIsList(t, output)
}
