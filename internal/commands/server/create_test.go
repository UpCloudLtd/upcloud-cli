package server

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/storage"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	Title1             = "mock-storage-title1"
	Title2             = "mock-storage-title2"
	UUID1              = "0127dfd6-3884-4079-a948-3a8881df1a7a"
	UUID3              = "012c61a6-b8f0-48c2-a63a-b4bf7d26a655"
	PrivateNetworkUUID = "03b5b0a0-ad4c-4817-9632-dafdb3ace5d9"
	MockPrivateIPv4    = "10.0.0.1"
	MockPrivateIPv6    = "fd42:42::1"
	MockPublicIPv4     = "192.0.2.0"
	MockPublicIPv6     = "2001:DB8::1"
)

func TestCreateServer(t *testing.T) {
	Plans := upcloud.Plans{
		Plans: []upcloud.Plan{
			{
				Name:        "1xCPU-2GB",
				StorageSize: 50,
			},
		},
	}
	Storage1 := upcloud.Storage{
		UUID:   UUID1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}

	StorageUbuntu2004 := upcloud.Storage{
		UUID:         "01000000-0000-4000-8000-000030200200",
		Title:        "Ubuntu Server 20.04 LTS (Focal Fossa)",
		Access:       "public",
		State:        "online",
		Type:         "template",
		Size:         4,
		TemplateType: "native",
	}

	StorageUbuntu2404 := upcloud.Storage{
		UUID:         "01000000-0000-4000-8000-000030240200",
		Title:        "Ubuntu Server 24.04 LTS (Noble Numbat)",
		Access:       "public",
		State:        "online",
		Type:         "template",
		Size:         4,
		TemplateType: "cloud-init",
	}

	Storage3 := upcloud.Storage{
		UUID:   UUID3,
		Title:  Title2,
		Access: "public",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   10,
		Tier:   "maxiops",
	}
	storages := &upcloud.Storages{
		Storages: []upcloud.Storage{
			Storage1,
			StorageUbuntu2004,
			StorageUbuntu2404,
			Storage3,
		},
	}
	serverDetailsMaint := upcloud.ServerDetails{
		Server: upcloud.Server{
			UUID:  UUID1,
			State: upcloud.ServerStateMaintenance,
		},
		VideoModel: "vga",
		Firewall:   "off",
	}

	serverDetailsStarted := serverDetailsMaint
	serverDetailsStarted.State = upcloud.ServerStateStarted
	sshKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIHicO0RGJyOSGeMtrmXK1upkkrL5yOrRdjNFl0FLwV00 Example public key"

	for _, test := range []struct {
		name            string
		args            []string
		createServerReq request.CreateServerRequest
		error           string
	}{
		{
			name: "use default values",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--ssh-keys", sshKey,
			},
			createServerReq: request.CreateServerRequest{
				VideoModel:       "vga",
				TimeZone:         "UTC",
				Plan:             "1xCPU-2GB",
				Hostname:         "example.com",
				Title:            "test-server",
				Zone:             "uk-lon1",
				Metadata:         upcloud.True,
				PasswordDelivery: "none",
				LoginUser: &request.LoginUser{
					CreatePassword: "no",
					SSHKeys:        []string{sshKey},
				},
				StorageDevices: request.CreateServerStorageDeviceSlice{request.CreateServerStorageDevice{
					Action:  "clone",
					Address: "virtio",
					Storage: StorageUbuntu2404.UUID,
					Title:   "example.com-OS",
					Size:    50,
					Type:    upcloud.StorageTypeDisk,
				}},
			},
		},
		{
			name: "use native template",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--os", "Ubuntu Server 20.04 LTS (Focal Fossa)",
				"--ssh-keys", sshKey,
			},
			createServerReq: request.CreateServerRequest{
				VideoModel:       "vga",
				TimeZone:         "UTC",
				Plan:             "1xCPU-2GB",
				Hostname:         "example.com",
				Title:            "test-server",
				Zone:             "uk-lon1",
				PasswordDelivery: "none",
				LoginUser: &request.LoginUser{
					CreatePassword: "no",
					SSHKeys:        []string{sshKey},
				},
				StorageDevices: request.CreateServerStorageDeviceSlice{request.CreateServerStorageDevice{
					Action:  "clone",
					Address: "virtio",
					Storage: StorageUbuntu2004.UUID,
					Title:   "example.com-OS",
					Size:    50,
					Type:    upcloud.StorageTypeDisk,
				}},
			},
		},
		{
			name: "server OS set, size larger than the minimum",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--os", Storage1.UUID,
				"--os-storage-size", "100",
				"--password-delivery", "email",
			},
			createServerReq: request.CreateServerRequest{
				VideoModel:       "vga",
				TimeZone:         "UTC",
				Plan:             "1xCPU-2GB",
				Hostname:         "example.com",
				Title:            "test-server",
				Zone:             "uk-lon1",
				PasswordDelivery: "email",
				LoginUser:        &request.LoginUser{CreatePassword: "yes"},
				StorageDevices: request.CreateServerStorageDeviceSlice{request.CreateServerStorageDevice{
					Action:  "clone",
					Address: "virtio",
					Storage: Storage1.UUID,
					Title:   "example.com-OS",
					Size:    100,
					Type:    upcloud.StorageTypeDisk,
				}},
			},
		},
		{
			name: "flags mapped to the correct field",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--avoid-host", "1234",
				"--host", "5678",
				"--boot-order", "cdrom,network",
				"--user-data", "example.com",
				"--cores", "12",
				"--memory", "4096",
				"--plan", "custom",
				"--password-delivery", "sms",
				"--server-group", "a4643646-8342-4324-4134-364138712378",
				"--simple-backup", "00,monthlies",
				"--time-zone", "EET",
				"--video-model", "VM",
				"--enable-firewall",
				"--enable-metadata",
				"--username", "johndoe",
				"--enable-remote-access",
				"--remote-access-type", upcloud.RemoteAccessTypeVNC,
				"--remote-access-password", "secret",
				"--label", "env=test",
			},
			createServerReq: request.CreateServerRequest{
				Hostname:             "example.com",
				Title:                "test-server",
				Zone:                 "uk-lon1",
				AvoidHost:            1234,
				Host:                 5678,
				BootOrder:            "cdrom,network",
				UserData:             "example.com",
				CoreNumber:           12,
				MemoryAmount:         4096,
				Plan:                 "custom",
				PasswordDelivery:     "sms",
				ServerGroup:          "a4643646-8342-4324-4134-364138712378",
				SimpleBackup:         "00,monthlies",
				TimeZone:             "EET",
				VideoModel:           "VM",
				Firewall:             "on",
				Metadata:             upcloud.True,
				RemoteAccessEnabled:  upcloud.True,
				RemoteAccessType:     upcloud.RemoteAccessTypeVNC,
				RemoteAccessPassword: "secret",
				LoginUser:            &request.LoginUser{CreatePassword: "yes", Username: "johndoe"},
				StorageDevices: request.CreateServerStorageDeviceSlice{request.CreateServerStorageDevice{
					Action:  "clone",
					Address: "virtio",
					Storage: StorageUbuntu2404.UUID,
					Title:   "example.com-OS",
					Size:    10,
					Type:    upcloud.StorageTypeDisk,
				}},
				Labels: &upcloud.LabelSlice{upcloud.Label{
					Key:   "env",
					Value: "test",
				}},
			},
		},
		{
			name: "multiple labels",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--password-delivery", "email",
				"--label", "env=test",
				"--label", "env1=test1",
				"--label", "only_key",
			},
			createServerReq: request.CreateServerRequest{
				VideoModel:       "vga",
				TimeZone:         "UTC",
				Plan:             "1xCPU-2GB",
				Hostname:         "example.com",
				Title:            "test-server",
				Zone:             "uk-lon1",
				Metadata:         upcloud.True,
				PasswordDelivery: "email",
				LoginUser:        &request.LoginUser{CreatePassword: "yes"},
				StorageDevices: request.CreateServerStorageDeviceSlice{request.CreateServerStorageDevice{
					Action:  "clone",
					Address: "virtio",
					Storage: StorageUbuntu2404.UUID,
					Title:   "example.com-OS",
					Size:    50,
					Type:    upcloud.StorageTypeDisk,
				}},
				Labels: &upcloud.LabelSlice{
					upcloud.Label{
						Key:   "env",
						Value: "test",
					},
					upcloud.Label{
						Key:   "env1",
						Value: "test1",
					},
					upcloud.Label{
						Key: "only_key",
					},
				},
			},
		},
		{
			name: "multiple storages",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--os-storage-encrypt",
				"--password-delivery", "email",
				"--storage", "action=create,address=virtio,encrypt=true,type=disk,size=20,title=new-storage",
				"--storage", fmt.Sprintf("action=clone,storage=%s,title=three-clone", Storage3.Title),
				"--storage", fmt.Sprintf("action=attach,storage=%s,type=cdrom", Storage1.Title),
			},
			createServerReq: request.CreateServerRequest{
				VideoModel:       "vga",
				TimeZone:         "UTC",
				Plan:             "1xCPU-2GB",
				Hostname:         "example.com",
				Title:            "test-server",
				Zone:             "uk-lon1",
				Metadata:         upcloud.True,
				PasswordDelivery: "email",
				LoginUser:        &request.LoginUser{CreatePassword: "yes"},
				StorageDevices: request.CreateServerStorageDeviceSlice{
					request.CreateServerStorageDevice{
						Action:    "clone",
						Address:   "virtio",
						Encrypted: upcloud.FromBool(true),
						Storage:   StorageUbuntu2404.UUID,
						Title:     "example.com-OS",
						Size:      50,
						Type:      upcloud.StorageTypeDisk,
					},
					request.CreateServerStorageDevice{
						Action:    "create",
						Address:   "virtio",
						Encrypted: upcloud.FromBool(true),
						Title:     "new-storage",
						Size:      20,
						Type:      upcloud.StorageTypeDisk,
					},
					request.CreateServerStorageDevice{
						Action:  "clone",
						Storage: Storage3.UUID,
						Title:   "three-clone",
					},
					request.CreateServerStorageDevice{
						Action:  "attach",
						Storage: Storage1.UUID,
						Type:    upcloud.StorageTypeCDROM,
					},
				},
			},
		},
		{
			name: "with networks",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--password-delivery", "email",
				"--network", "type=public",
				"--network", "family=IPv4,type=public",
				"--network", fmt.Sprintf("type=public,ip-address=%s", MockPublicIPv4),
				"--network", fmt.Sprintf("type=public,ip-address=%s", MockPublicIPv6),
				"--network", "family=IPv6,type=public",
				"--network", "family=IPv4,type=utility",
				"--network", fmt.Sprintf("family=IPv4,type=private,network=%s,enable-bootable,disable-source-ip-filtering", PrivateNetworkUUID),
				"--network", fmt.Sprintf("type=private,network=%s,ip-address=%s", PrivateNetworkUUID, MockPrivateIPv4),
				"--network", fmt.Sprintf("family=IPv6,type=private,network=%s", PrivateNetworkUUID),
				"--network", fmt.Sprintf("type=private,network=%s,ip-address=%s", PrivateNetworkUUID, MockPrivateIPv6),
			},
			createServerReq: request.CreateServerRequest{
				VideoModel:       "vga",
				TimeZone:         "UTC",
				Plan:             "1xCPU-2GB",
				Hostname:         "example.com",
				Title:            "test-server",
				Zone:             "uk-lon1",
				Metadata:         upcloud.True,
				PasswordDelivery: "email",
				LoginUser:        &request.LoginUser{CreatePassword: "yes"},
				StorageDevices: request.CreateServerStorageDeviceSlice{request.CreateServerStorageDevice{
					Action:  "clone",
					Address: "virtio",
					Storage: StorageUbuntu2404.UUID,
					Title:   "example.com-OS",
					Size:    50,
					Type:    upcloud.StorageTypeDisk,
				}},
				Networking: &request.CreateServerNetworking{Interfaces: request.CreateServerInterfaceSlice{
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: upcloud.IPAddressFamilyIPv4}},
						Type:        upcloud.NetworkTypePublic,
					},
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: upcloud.IPAddressFamilyIPv4}},
						Type:        upcloud.NetworkTypePublic,
					},
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: upcloud.IPAddressFamilyIPv4, Address: MockPublicIPv4}},
						Type:        upcloud.NetworkTypePublic,
					},
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: upcloud.IPAddressFamilyIPv6, Address: MockPublicIPv6}},
						Type:        upcloud.NetworkTypePublic,
					},
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: upcloud.IPAddressFamilyIPv6}},
						Type:        upcloud.NetworkTypePublic,
					},
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: upcloud.IPAddressFamilyIPv4}},
						Type:        upcloud.NetworkTypeUtility,
					},
					request.CreateServerInterface{
						IPAddresses:       request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: upcloud.IPAddressFamilyIPv4}},
						Type:              upcloud.NetworkTypePrivate,
						Network:           PrivateNetworkUUID,
						Bootable:          upcloud.True,
						SourceIPFiltering: upcloud.False,
					},
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: upcloud.IPAddressFamilyIPv4, Address: MockPrivateIPv4}},
						Type:        upcloud.NetworkTypePrivate,
						Network:     PrivateNetworkUUID,
					},
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: upcloud.IPAddressFamilyIPv6}},
						Type:        upcloud.NetworkTypePrivate,
						Network:     PrivateNetworkUUID,
					},
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: upcloud.IPAddressFamilyIPv6, Address: MockPrivateIPv6}},
						Type:        upcloud.NetworkTypePrivate,
						Network:     PrivateNetworkUUID,
					},
				}},
			},
		},
		{
			name: "networks type missing",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--network", "family=IPv4,type=utility",
				"--network", "family=IPv6,type=public",
				"--network", "family=IPv6",
				"--password-delivery", "sms",
			},
			error: "network type is required",
		},
		{
			name: "invalid ip address",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--network", "type=public,ip-address=10.0.0.300",
				"--password-delivery", "sms",
			},
			error: "10.0.0.300 is an invalid ip address",
		},
		{
			name: "hostname is missing",
			args: []string{
				"--title", "title",
				"--zone", "zone",
			},
			error: `required flag(s) "hostname" not set`,
		},
		{
			name: "zone is missing",
			args: []string{
				"--title", "title",
				"--hostname", "hostname",
			},
			error: `required flag(s) "zone" not set`,
		},
		{
			name: "label argument missing",
			args: []string{
				"--title", "title",
				"--hostname", "hostname",
				"--zone", "zone",
				"--label", "env=test",
				"--label",
			},
			error: `flag needs an argument: --label`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := CreateCommand()
			mService := new(smock.Service)

			storage.CachedStorages = nil
			createServerReq := test.createServerReq
			mService.On("CreateServer", &createServerReq).Return(&serverDetailsMaint, nil)
			mService.On("GetPlans", mock.Anything).Return(&Plans, nil)
			mService.On("GetStorages", mock.Anything).Return(storages, nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.error != "" {
				if err == nil {
					t.Errorf("expected error '%v', got nil", test.error)
				} else {
					assert.Equal(t, test.error, err.Error())
				}
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, "GetStorages", 1)
				mService.AssertNumberOfCalls(t, "CreateServer", 1)
			}
		})
	}
}
