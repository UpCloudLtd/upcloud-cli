package kubernetes

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCommand(t *testing.T) {
	targetMethod := "DeleteKubernetesCluster"

	cluster := upcloud.KubernetesCluster{
		UUID: "17fbd082-30b0-11eb-adc1-0242ac120003",
		Name: "test-cluster",
	}

	for _, test := range []struct {
		name  string
		arg   string
		error string
		req   request.DeleteKubernetesClusterRequest
	}{
		{
			name: "delete with UUID",
			arg:  cluster.UUID,
			req:  request.DeleteKubernetesClusterRequest{UUID: cluster.UUID},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			req := test.req
			mService.On(targetMethod, &req).Return(nil)

			conf := config.New()
			c := commands.BuildCommand(DeleteCommand(), nil, conf)

			_, err := c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, &mService, flume.New("test")), test.arg)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
