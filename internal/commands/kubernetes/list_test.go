package kubernetes

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/testutils"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
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
