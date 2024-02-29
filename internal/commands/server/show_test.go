package server

import (
	"context"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

func TestServerHumanOutput(t *testing.T) {
	text.DisableColors()
	uuid := "0077fa3d-32db-4b09-9f5f-30d9e9afb565"
	srv := &upcloud.ServerDetails{
		Server: upcloud.Server{
			CoreNumber:   0,
			Hostname:     "server1.example.com",
			License:      0,
			MemoryAmount: 2048,
			State:        "started",
			Plan:         "1xCPU-2GB",
			Title:        "server1.example.com",
			UUID:         uuid,
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
		ServerGroup:  "a4643646-8342-4324-4134-364138712378",
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
				Direction: upcloud.FirewallRuleDirectionIn,
				Action:    upcloud.FirewallRuleActionAccept,
				Family:    upcloud.IPAddressFamilyIPv4,
				Protocol:  upcloud.FirewallRuleProtocolTCP,
				Position:  1,
				Comment:   "This is the comment",
			},
			{
				Direction:               upcloud.FirewallRuleDirectionOut,
				Action:                  upcloud.FirewallRuleActionDrop,
				Family:                  upcloud.IPAddressFamilyIPv4,
				Protocol:                upcloud.FirewallRuleProtocolUDP,
				ICMPType:                upcloud.FirewallRuleProtocolICMP,
				Position:                2,
				SourceAddressStart:      "10.10.10.0",
				SourceAddressEnd:        "10.10.10.99",
				DestinationAddressStart: "10.20.20.0",
				SourcePortStart:         "0",
				SourcePortEnd:           "1024",
			},
		},
	}

	expected := `  
  Common
    UUID:          0077fa3d-32db-4b09-9f5f-30d9e9afb565 
    Hostname:      server1.example.com                  
    Title:         server1.example.com                  
    Plan:          1xCPU-2GB                            
    Zone:          fi-hel1                              
    State:         started                              
    Simple Backup: 0100,dailies                         
    Licence:       0                                    
    Metadata:      True                                 
    Timezone:      UTC                                  
    Host ID:       7653311107                           
    Server Group:  a4643646-8342-4324-4134-364138712378 
    Tags:          DEV,Ubuntu                           

  Labels:

    No labels defined for this resource.
    
  Storage: (Flags: B = bootdisk, P = part of plan)

     UUID                                   Title                             Type   Address    Size (GiB)   Encrypted   Flags 
    ────────────────────────────────────── ───────────────────────────────── ────── ────────── ──────────── ─────────── ───────
     012580a1-32a1-466e-a323-689ca16f2d43   Storage for server1.example.com   disk   virtio:0           20   no          P     
    
  NICs: (Flags: S = source IP filtering, B = bootable)

     #   Type      IP Address                                           MAC Address         Network                                Flags 
    ─── ───────── ──────────────────────────────────────────────────── ─────────────────── ────────────────────────────────────── ───────
     1   public    IPv4: 94.237.0.207                                   de:ff:ff:ff:66:89   037fcf2a-6745-45dd-867e-f9479ea8c044         
     2   utility   IPv4: 10.6.3.95 (f)                                  de:ff:ff:ff:ed:85   03000000-0000-4000-8045-000000000000         
     3   public    IPv6: xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx   de:ff:ff:ff:cc:20   03c93fd8-cc60-4849-91b8-6e404b228e2a         
    
  Firewall Rules:

     #   Direction   Action   Src IPAddress   Dest IPAddress   Src Port   Dest Port   Protocol      
    ─── ─────────── ──────── ─────────────── ──────────────── ────────── ─────────── ───────────────
     1   in          accept   Any             Any              Any        Any         IPv4/tcp      
     2   out         drop     10.10.10.0 →    10.20.20.0       0 →        Any         IPv4/udp/icmp 
                              10.10.10.99                      1024                                 
    
  
  Remote Access
    Type:     vnc                     
    Host:     fi-hel1.vnc.upcloud.com 
    Port:     3000                    
    Password: aabbccdd                

`

	mService := smock.Service{}
	mService.On("GetServers").Return(&upcloud.Servers{Servers: []upcloud.Server{srv.Server}}, nil)
	mService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: uuid}).Return(srv, nil)
	mService.On("GetFirewallRules", &request.GetFirewallRulesRequest{ServerUUID: uuid}).Return(firewallRules, nil)

	conf := config.New()

	command := commands.BuildCommand(ShowCommand(), nil, conf)

	// get resolver to initialize command cache
	_, err := command.(*showCommand).Get(context.TODO(), &mService)
	if err != nil {
		t.Fatal(err)
	}

	command.Cobra().SetArgs([]string{uuid})
	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}
