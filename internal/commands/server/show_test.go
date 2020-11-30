package server

import (
	"bytes"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

func TestServerHumanOutput(t *testing.T) {
	text.DisableColors()
	s := &upcloud.ServerDetails{
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
							Floating: upcloud.False,
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
				Direction: upcloud.FirewallRuleDirectionIn,
				Action:    upcloud.FirewallRuleActionAccept,
				Family:    upcloud.IPAddressFamilyIPv4,
				Protocol:  upcloud.FirewallRuleProtocolTCP,
				Position:  1,
				Comment:   "This is the comment",
			},
		},
	}

	expected := `  
  Common:
    UUID:     0077fa3d-32db-4b09-9f5f-30d9e9afb565 
    Title:    server1.example.com                  
    Hostname: server1.example.com                  
    Plan:     1xCPU-2GB                            
    Zone:     fi-hel1                              
    State:    started                              
    Tags:     DEV,Ubuntu                           
    Licence:  0                                    
    Metadata: yes                                  
    Timezone: UTC                                  
    Host ID:  7653311107                           
  
  Storage:
    Simple Backup: 0100,dailies 
    Devices:
      
       Title (UUID)                             Type   Address    Size (GiB)   Flags 
      ──────────────────────────────────────── ────── ────────── ──────────── ───────
       Storage for server1.example.com          disk   virtio:0           20   P     
       (012580a1-32a1-466e-a323-689ca16f2d43)                                        
      
      (Flags: B = bootdisk, P = part of plan)
  
  Networking:
    NICS:
      
       #   Type      Network               Addresses                                            Flags 
      ─── ───────── ───────────────────── ──────────────────────────────────────────────────── ───────
       1   public    037fcf2a-6745-45dd-   MAC:  de:ff:ff:ff:66:89                                    
                     867e-f9479ea8c044     IPv4: 94.237.0.207                                         
       2   utility   03000000-0000-4000-   MAC:  de:ff:ff:ff:ed:85                                    
                     8045-000000000000     IPv4: 10.6.3.95                                            
       3   public    03c93fd8-cc60-4849-   MAC:  de:ff:ff:ff:cc:20                                    
                     91b8-6e404b228e2a     IPv6: xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx         
      
      (Flags: S = source IP filtering, B = bootable)
  
  Firewall Rules:
    
     #   Action   Source   Destination   Dir   Proto    
    ─── ──────── ──────── ───────────── ───── ──────────
     1   accept                          in    IPv4/tcp 
  
  Remote Access:
    Enabled   yes                          
    Type:     vnc                          
    Address:  fi-hel1.vnc.upcloud.com:3000 
    Password: aabbccdd                     
`

	buf := new(bytes.Buffer)
	command := ShowCommand(&MockServerService{}, &MockFirewallService{})
	err := command.HandleOutput(buf, &commandResponseHolder{s, firewallRules})

	assert.Nil(t, err)
	assert.Equal(t, expected, buf.String())
}
