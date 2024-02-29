package servergroup

import (
	"context"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

func TestShowCommand(t *testing.T) {
	text.DisableColors()

	server1 := upcloud.ServerDetails{
		Server: upcloud.Server{
			Hostname: "test1",
			State:    upcloud.ServerStateStarted,
			Title:    "test1",
			UUID:     "00fda506-2e8f-44f7-9d06-74d197d029f3",
			Zone:     "pl-waw1",
		},
		Host:        4815602307,
		ServerGroup: "0bb022ba-7ab5-4c94-8671-2e32fa543a79",
	}

	server2 := upcloud.ServerDetails{
		Server: upcloud.Server{
			Hostname: "test2",
			State:    upcloud.ServerStateStarted,
			Title:    "test2",
			UUID:     "00333d1b-3a4a-4b75-820a-4a56d70395dd",
			Zone:     "pl-waw1",
		},
		Host:        8573698691,
		ServerGroup: "0bb022ba-7ab5-4c94-8671-2e32fa543a79",
	}

	server3 := upcloud.ServerDetails{
		Server: upcloud.Server{
			Hostname: "test3",
			State:    upcloud.ServerStateStarted,
			Title:    "test3",
			UUID:     "33333d1b-3a4a-4b75-820a-4a56d7039533",
			Zone:     "pl-waw1",
		},
		Host:        8573698691,
		ServerGroup: "0bb022ba-7ab5-4c94-8671-2e32fa543a79",
	}

	serverGroup1 := upcloud.ServerGroup{
		AntiAffinityPolicy: upcloud.ServerGroupAntiAffinityPolicyStrict,
		AntiAffinityStatus: []upcloud.ServerGroupMemberAntiAffinityStatus{
			{
				ServerUUID: server1.UUID,
				Status:     "met",
			},
			{
				ServerUUID: server2.UUID,
				Status:     "unmet",
			},
			{
				ServerUUID: server3.UUID,
				Status:     "unmet",
			},
		},
		Labels: []upcloud.Label{
			{
				Key:   "managedBy",
				Value: "upcloud-cli-unit-test",
			},
			{
				Key:   "another",
				Value: "label-thing",
			},
		},
		Members: upcloud.ServerUUIDSlice{
			server1.UUID,
			server2.UUID,
			server3.UUID,
		},
		Title: "test1",
		UUID:  "0bb022ba-7ab5-4c94-8671-2e32fa543a79",
	}

	expected := `  
  Overview:
    UUID:                 0bb022ba-7ab5-4c94-8671-2e32fa543a79 
    Title:                test1                                
    Anti-affinity policy: strict                               
    Anti-affinity state:  unmet                                
    Server count:         3                                    

  Labels:

     Key         Value                 
    ─────────── ───────────────────────
     managedBy   upcloud-cli-unit-test 
     another     label-thing           
    
  Servers:

     UUID                                   Hostname   Zone      Host         Anti-affinity state   State   
    ────────────────────────────────────── ────────── ───────── ──────────── ───────────────────── ─────────
     00fda506-2e8f-44f7-9d06-74d197d029f3   test1      pl-waw1   4815602307   met                   started 
     00333d1b-3a4a-4b75-820a-4a56d70395dd   test2      pl-waw1   8573698691   unmet                 started 
     33333d1b-3a4a-4b75-820a-4a56d7039533   test3      pl-waw1   8573698691   unmet                 started 
    
`

	mService := smock.Service{}
	mService.On("GetServerGroups", &request.GetServerGroupsRequest{}).Return(upcloud.ServerGroups{serverGroup1}, nil)
	mService.On("GetServerGroup", &request.GetServerGroupRequest{UUID: serverGroup1.UUID}).Return(&serverGroup1, nil)
	mService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: server1.UUID}).Return(&server1, nil)
	mService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: server2.UUID}).Return(&server2, nil)
	mService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: server3.UUID}).Return(&server3, nil)

	conf := config.New()
	command := commands.BuildCommand(ShowCommand(), nil, conf)

	// get resolver to initialize command cache
	_, err := command.(*showCommand).Get(context.TODO(), &mService)
	if err != nil {
		t.Fatal(err)
	}

	command.Cobra().SetArgs([]string{serverGroup1.UUID})
	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}
