package token

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCommand(t *testing.T) {
	svc := &smock.Service{}
	conf := config.New()
	dr := &request.DeleteTokenRequest{ID: "0cdabbf9-090b-4fc5-a6ae-3f76801ed171"}

	svc.On("DeleteToken", dr).Once().Return(nil)

	command := commands.BuildCommand(DeleteCommand(), nil, conf)
	_, err := command.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, svc, flume.New("test")), dr.ID)
	assert.NoError(t, err)

	svc.AssertExpectations(t)
}
