package zone

import (
	"bytes"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/gemalto/flume"
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
	// force human output
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(ListCommand(), nil, conf)

	res, err := command.(commands.NoArgumentCommand).ExecuteWithoutArguments(commands.NewExecutor(conf, &mService, flume.New("test")))

	assert.Nil(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, conf, res)
	assert.NoError(t, err)
	assert.Regexp(t, "ID\\s+Description\\s+Public", buf.String())
	assert.Regexp(t, "fi-hel1\\s+Helsinki #1\\s+yes", buf.String())
}
