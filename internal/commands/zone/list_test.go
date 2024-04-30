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

	for _, test := range []struct {
		name     string
		zones    upcloud.Zones
		expected string
	}{
		{
			name: "no private zones",
			zones: upcloud.Zones{
				Zones: []upcloud.Zone{
					{ID: "fi-hel1", Description: "Helsinki #1", Public: 1},
					{ID: "de-fra1", Description: "Frankfurt #1", Public: 1},
				},
			},
			expected: `
 ID        Description    Public 
───────── ────────────── ────────
 fi-hel1   Helsinki #1    yes    
 de-fra1   Frankfurt #1   yes    

`,
		}, {
			name: "with private zones",
			zones: upcloud.Zones{
				Zones: []upcloud.Zone{
					{ID: "de-fra1", Description: "Frankfurt #1", Public: 1},
					{ID: "de-tst1", Description: "Test #1", Public: 0, ParentZone: "de-fra1"},
				},
			},
			expected: `
 ID        Description    Public   Parent zone 
───────── ────────────── ──────── ─────────────
 de-fra1   Frankfurt #1   yes                  
 de-tst1   Test #1        no       de-fra1     

`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			zones := test.zones
			mService := smock.Service{}
			mService.On("GetZones").Return(&zones, nil)

			conf := config.New()
			command := commands.BuildCommand(ListCommand(), nil, conf)

			output, err := mockexecute.MockExecute(command, &mService, conf)

			assert.NoError(t, err)
			assert.Equal(t, output, test.expected)
		})
	}
}
