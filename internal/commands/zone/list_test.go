package zone

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

func TestZoneListHumanOutput(t *testing.T) {
	text.DisableColors()
	zones := upcloud.Zones{
		Zones: []upcloud.Zone{
			{ID: "fi-hel1", Description: "Helsinki #1", Public: 1},
			{ID: "de-fra1", Description: "Frankfurt #1", Public: 1},
		},
	}

	mService := smock.Service{}
	mService.On("GetZones").Return(&zones, nil)

	conf := config.New()
	command := commands.BuildCommand(ListCommand(), nil, conf)

	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.NoError(t, err)
	assert.Regexp(t, "ID\\s+Description\\s+Public", output)
	assert.Regexp(t, "fi-hel1\\s+Helsinki #1\\s+yes", output)
}
