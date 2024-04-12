package servergroup

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestModifyCommand(t *testing.T) {
	server1 := upcloud.ServerDetails{
		Server: upcloud.Server{
			Hostname: "test1",
			Title:    "test1",
			UUID:     "11111111-1111-1111-1111-111111111111",
			Zone:     "pl-waw1",
		},
		Host:        4815602307,
		ServerGroup: "8abc8009-4325-4b23-4321-b1232cd81231",
	}

	server2 := upcloud.ServerDetails{
		Server: upcloud.Server{
			Hostname: "test2",
			Title:    "test2",
			UUID:     "22222222-2222-2222-2222-222222222222",
			Zone:     "pl-waw1",
		},
		Host:        8573698691,
		ServerGroup: "8abc8009-4325-4b23-4321-b1232cd81231",
	}

	serverGroup := upcloud.ServerGroup{
		AntiAffinityPolicy: upcloud.ServerGroupAntiAffinityPolicyBestEffort,
		Labels: []upcloud.Label{
			{
				Key:   "env",
				Value: "original",
			},
		},
		Members: []string{
			server1.UUID,
			server2.UUID,
		},
		Title: "test-server-group",
		UUID:  "8abc8009-4325-4b23-4321-b1232cd81231",
	}

	for _, test := range []struct {
		name    string
		args    []string
		error   string
		returns *upcloud.ServerGroup
		req     request.ModifyServerGroupRequest
	}{
		{
			name: "title is passed",
			args: []string{"--title", "New title", serverGroup.UUID},
			returns: &upcloud.ServerGroup{
				AntiAffinityPolicy: serverGroup.AntiAffinityPolicy,
				Labels:             serverGroup.Labels,
				Members:            serverGroup.Members,
				Title:              "New title",
				UUID:               serverGroup.UUID,
			},
			req: request.ModifyServerGroupRequest{
				Title: "New title",
				UUID:  serverGroup.UUID,
			},
		},
		{
			name: "servers are passed",
			args: []string{
				serverGroup.UUID,
				"--server", "11111111-1111-1111-1111-111111111111",
				"--server", "22222222-2222-2222-2222-222222222222",
			},
			returns: &upcloud.ServerGroup{
				AntiAffinityPolicy: serverGroup.AntiAffinityPolicy,
				Labels:             serverGroup.Labels,
				Members: []string{
					"11111111-1111-1111-1111-111111111111",
					"22222222-2222-2222-2222-222222222222",
				},
				Title: serverGroup.Title,
				UUID:  serverGroup.UUID,
			},
			req: request.ModifyServerGroupRequest{
				Members: &upcloud.ServerUUIDSlice{
					"11111111-1111-1111-1111-111111111111",
					"22222222-2222-2222-2222-222222222222",
				},
				UUID: serverGroup.UUID,
			},
		},
		{
			name: "labels are passed",
			args: []string{
				serverGroup.UUID,
				"--label", "env=test",
				"--label", "another=label",
			},
			returns: &upcloud.ServerGroup{
				AntiAffinityPolicy: serverGroup.AntiAffinityPolicy,
				Labels: []upcloud.Label{
					{
						Key:   "env",
						Value: "test",
					},
					{
						Key:   "another",
						Value: "label",
					},
				},
				Members: serverGroup.Members,
				Title:   serverGroup.Title,
				UUID:    serverGroup.UUID,
			},
			req: request.ModifyServerGroupRequest{
				Labels: &upcloud.LabelSlice{
					{
						Key:   "env",
						Value: "test",
					},
					{
						Key:   "another",
						Value: "label",
					},
				},
				UUID: serverGroup.UUID,
			},
		},
		{
			name: "anti-affinity-policy and members are passed",
			args: []string{
				serverGroup.UUID,
				"--anti-affinity-policy", "strict",
				"--server", "11111111-1111-1111-1111-111111111111",
				"--server", "22222222-2222-2222-2222-222222222222",
			},
			returns: &upcloud.ServerGroup{
				AntiAffinityPolicy: upcloud.ServerGroupAntiAffinityPolicyStrict,
				Labels:             serverGroup.Labels,
				Members: []string{
					"11111111-1111-1111-1111-111111111111",
					"22222222-2222-2222-2222-222222222222",
				},
				Title: serverGroup.Title,
				UUID:  serverGroup.UUID,
			},
			req: request.ModifyServerGroupRequest{
				AntiAffinityPolicy: upcloud.ServerGroupAntiAffinityPolicyStrict,
				Members: &upcloud.ServerUUIDSlice{
					"11111111-1111-1111-1111-111111111111",
					"22222222-2222-2222-2222-222222222222",
				},
				UUID: serverGroup.UUID,
			},
		},
	} {
		targetMethod := "ModifyServerGroup"
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			mService.On("GetServers", mock.Anything).Return(
				&upcloud.Servers{
					Servers: []upcloud.Server{
						server1.Server,
						server2.Server,
					},
				},
				nil,
			)
			req := test.req
			mService.On(targetMethod, &req).Return(test.returns, nil)
			mService.On("GetServerGroups", mock.Anything).Return(upcloud.ServerGroups{serverGroup}, nil)

			conf := config.New()

			c := commands.BuildCommand(ModifyCommand(), nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, &mService, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
