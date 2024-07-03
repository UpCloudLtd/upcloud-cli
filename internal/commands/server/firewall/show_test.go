package serverfirewall

import (
	"bytes"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

func TestFirewallShowHumanOutput(t *testing.T) {
	text.DisableColors()
	serverDetails := &upcloud.ServerDetails{
		Server: upcloud.Server{
			CoreNumber:   0,
			Hostname:     "server1.example.com",
			License:      0,
			MemoryAmount: 2048,
			State:        "started",
			Plan:         "1xCPU-2GB",
			Title:        "server1.example.com",
			UUID:         "0077fa3d-32db-4b09-9f5f-30d9e9afb565",
			Zone:         "fi-hel1",
			Tags: []string{
				"DEV",
				"Ubuntu",
			},
		},
		BootOrder: "cdrom,disk",
		Firewall:  "on",
		Host:      7653311107,
		IPAddresses: []upcloud.IPAddress{
			{
				Access:  "private",
				Address: "10.0.0.00",
				Family:  "IPv4",
			},
			{
				Access:  "public",
				Address: "0.0.0.0",
				Family:  "IPv4",
			},
			{
				Access:  "public",
				Address: "xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx",
				Family:  "IPv6",
			},
		},
		Metadata: upcloud.True,
		Networking: upcloud.ServerNetworking{
			Interfaces: []upcloud.ServerInterface{
				{
					Index: 1,
					IPAddresses: []upcloud.IPAddress{
						{
							Address:  "94.237.0.207",
							Family:   "IPv4",
							Floating: upcloud.False,
						},
					},
					MAC:      "de:ff:ff:ff:66:89",
					Network:  "037fcf2a-6745-45dd-867e-f9479ea8c044",
					Type:     "public",
					Bootable: upcloud.False,
				},
				{
					Index: 2,
					IPAddresses: []upcloud.IPAddress{
						{
							Address:  "10.6.3.95",
							Family:   "IPv4",
							Floating: upcloud.True,
						},
					},
					MAC:      "de:ff:ff:ff:ed:85",
					Network:  "03000000-0000-4000-8045-000000000000",
					Type:     "utility",
					Bootable: upcloud.False,
				},
				{
					Index: 3,
					IPAddresses: []upcloud.IPAddress{
						{
							Address:  "xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx",
							Family:   "IPv6",
							Floating: upcloud.False,
						},
					},
					MAC:      "de:ff:ff:ff:cc:20",
					Network:  "03c93fd8-cc60-4849-91b8-6e404b228e2a",
					Type:     "public",
					Bootable: upcloud.False,
				},
			},
		},
		NICModel:     "virtio",
		SimpleBackup: "0100,dailies",
		StorageDevices: []upcloud.ServerStorageDevice{
			{
				Address:    "virtio:0",
				PartOfPlan: "yes",
				UUID:       "012580a1-32a1-466e-a323-689ca16f2d43",
				Size:       20,
				Title:      "Storage for server1.example.com",
				Type:       "disk",
				BootDisk:   0,
			},
		},
		Timezone:             "UTC",
		VideoModel:           "cirrus",
		RemoteAccessEnabled:  upcloud.True,
		RemoteAccessType:     "vnc",
		RemoteAccessHost:     "fi-hel1.vnc.upcloud.com",
		RemoteAccessPassword: "aabbccdd",
		RemoteAccessPort:     3000,
	}

	firewallRules := &upcloud.FirewallRules{
		FirewallRules: []upcloud.FirewallRule{
			{
				SourceAddressStart:      "192.168.1.1",
				SourcePortStart:         "10",
				SourcePortEnd:           "20",
				DestinationAddressStart: "127.0.0.1",
				DestinationAddressEnd:   "127.0.0.127",
				DestinationPortStart:    "30",
				DestinationPortEnd:      "30",
				Direction:               upcloud.FirewallRuleDirectionIn,
				Action:                  upcloud.FirewallRuleActionAccept,
				Family:                  upcloud.IPAddressFamilyIPv4,
				Protocol:                upcloud.FirewallRuleProtocolTCP,
				Position:                1,
				Comment:                 "This is the comment",
			},
		},
	}

	expected := `
  Firewall rules

     #   Action   Source          Destination   Dir   Proto    
    ─── ──────── ─────────────── ───────────── ───── ──────────
     1   accept   192.168.1.1 →   127.0.0.1 →   in    IPv4/tcp 
                  <nil>           127.0.0.127                  
                  port: 10 → 20   port: 30                     
    
  
  Enabled yes 

`

	conf := config.New()
	testCmd := ShowCommand()
	mService := new(smock.Service)

	mService.On("GetFirewallRules",
		&request.GetFirewallRulesRequest{ServerUUID: serverDetails.UUID},
	).Return(firewallRules, nil)
	mService.On("GetServerDetails",
		&request.GetServerDetailsRequest{UUID: serverDetails.UUID},
	).Return(serverDetails, nil)
	// force human output
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(testCmd, nil, conf)
	out, err := command.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, mService, flume.New("test")), serverDetails.UUID)
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, conf.Output(), out)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())

	mService.AssertNumberOfCalls(t, "GetFirewallRules", 1)
	mService.AssertNumberOfCalls(t, "GetServerDetails", 1)
}
