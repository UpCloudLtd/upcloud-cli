package server

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestModifyCommand(t *testing.T) {
	targetMethod := "ModifyServer"

	Server1 := upcloud.Server{
		CoreNumber:   1,
		Hostname:     "server-1-hostname",
		License:      0,
		MemoryAmount: 1024,
		Plan:         "server-1-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-1-title",
		UUID:         "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		Zone:         "fi-hel1",
	}

	details := upcloud.ServerDetails{
		Server: Server1,
	}

	for _, test := range []struct {
		name       string
		args       []string
		server     upcloud.Server
		modifyCall request.ModifyServerRequest
	}{
		{
			name: "Backend called, flags mapped to the correct field",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--boot-order", "cdrom,network",
				"--cores", "12",
				"--memory", "4096",
				"--plan", "custom",
				"--simple-backup", "00,monthlies",
				"--time-zone", "EET",
				"--video-model", "VM",
				"--enable-firewall",
				"--enable-metadata",
				"--enable-remote-access",
				"--remote-access-type", upcloud.RemoteAccessTypeVNC,
				"--remote-access-password", "secret",
				"--label", "env=test",
			},
			server: Server1,
			modifyCall: request.ModifyServerRequest{
				UUID:                 Server1.UUID,
				Hostname:             "example.com",
				Title:                "test-server",
				BootOrder:            "cdrom,network",
				CoreNumber:           12,
				MemoryAmount:         4096,
				Plan:                 "custom",
				SimpleBackup:         "00,monthlies",
				TimeZone:             "EET",
				VideoModel:           "VM",
				Firewall:             "on",
				Metadata:             upcloud.FromBool(true),
				RemoteAccessEnabled:  upcloud.FromBool(true),
				RemoteAccessType:     upcloud.RemoteAccessTypeVNC,
				RemoteAccessPassword: "secret",
				Labels: &upcloud.LabelSlice{upcloud.Label{
					Key:   "env",
					Value: "test",
				}},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := ModifyCommand()
			mService := new(smock.Service)

			modifyCall := test.modifyCall
			mService.On(targetMethod, &modifyCall).Return(&details, nil)
			mService.On("GetServerDetails", mock.Anything).Return(&details, nil)
			c := commands.BuildCommand(testCmd, nil, conf)
			err := c.Cobra().Flags().Parse(test.args)
			assert.NoError(t, err)

			_, err = c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, mService, flume.New("test")), test.server.UUID)
			assert.NoError(t, err)
			mService.AssertNumberOfCalls(t, targetMethod, 1)
			mService.AssertNumberOfCalls(t, "GetServerDetails", 1)
		})
	}
}
