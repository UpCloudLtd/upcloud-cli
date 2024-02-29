package kubernetes

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var versions = []upcloud.KubernetesVersion{
	{
		Id:      "1.26",
		Version: "v1.26.3",
	},
	{
		Id:      "1.27",
		Version: "v1.27.4",
	},
}

func TestVersionsCommand(t *testing.T) {
	text.DisableColors()

	expected := `
 ID     Version 
────── ─────────
 1.26   v1.26.3 
 1.27   v1.27.4 

`

	mService := smock.Service{}
	mService.On("GetKubernetesVersions", mock.Anything).Return(versions, nil)

	conf := config.New()
	command := commands.BuildCommand(VersionsCommand(), nil, conf)

	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}
