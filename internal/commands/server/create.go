package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/ipaddress"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/storage"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
)

const defaultIPAddressFamily = upcloud.IPAddressFamilyIPv4

// CreateCommand creates the "server create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a server",
			"upctl server create --title myapp --zone fi-hel1 --hostname myapp --password-delivery email",
			"upctl server create --title \"My Server\" --zone fi-hel1 --hostname myapp --password-delivery email",
			"upctl server create --zone fi-hel1 --hostname myapp --password-delivery email --plan 2xCPU-4GB",
			"upctl server create --zone fi-hel1 --hostname myapp --password-delivery email --os \"Debian GNU/Linux 10 (Buster)\"",
			"upctl server create --zone fi-hel1 --hostname myapp --ssh-keys /path/to/publickey --network type=private,network=037a530b-533e-4cef-b6ad-6af8094bb2bc,ip-address=10.0.0.1",
		),
	}
}

var defaultCreateParams = &createParams{
	CreateServerRequest: request.CreateServerRequest{
		VideoModel:       "vga",
		TimeZone:         "UTC",
		Plan:             "1xCPU-2GB",
		PasswordDelivery: request.PasswordDeliveryNone,
	},
	firewall:       false,
	metadata:       false,
	os:             "Ubuntu Server 20.04 LTS (Focal Fossa)",
	osStorageSize:  0,
	sshKeys:        nil,
	username:       "",
	createPassword: false,
}

type createParams struct {
	request.CreateServerRequest
	firewall      bool
	metadata      bool
	os            string
	osStorageSize int

	storages []string
	networks []string

	sshKeys        []string
	username       string
	createPassword bool
	remoteAccess   bool
}

func (s *createParams) processParams(planSvc service.Plans, storageSvc service.Storage) error {
	if s.os != "" {
		var osStorage *upcloud.Storage

		osStorage, err := storage.SearchSingleStorage(s.os, storageSvc)
		if err != nil {
			return err
		}

		plans, err := planSvc.GetPlans()
		if err != nil {
			return err
		}

		size := minStorageSize
		if s.Plan != customPlan {
			for _, plan := range plans.Plans {
				if plan.Name == s.Plan {
					size = plan.StorageSize
				}
			}
		}

		s.StorageDevices = append(s.StorageDevices, request.CreateServerStorageDevice{
			Action:  "clone",
			Storage: osStorage.UUID,
			Title:   fmt.Sprintf("%s-OS", ui.TruncateText(s.Hostname, 64-7)),
			Size:    size,
			Tier:    upcloud.StorageTierMaxIOPS,
			Type:    upcloud.StorageTypeDisk,
		})
	}

	if s.osStorageSize != 0 {
		s.StorageDevices[0].Size = s.osStorageSize
	}

	if s.LoginUser == nil {
		s.LoginUser = &request.LoginUser{}
	}
	s.LoginUser.CreatePassword = "no"
	if s.username != "" {
		s.LoginUser.Username = s.username
	}
	return nil
}

func (s *createParams) handleStorage(in string, storageSvc service.Storage) (*request.CreateServerStorageDevice, error) {
	sd := &request.CreateServerStorageDevice{}
	fs := &pflag.FlagSet{}
	args, err := commands.Parse(in)
	if err != nil {
		return nil, err
	}
	fs.StringVar(&sd.Action, "action", sd.Action, "")
	fs.StringVar(&sd.Address, "address", sd.Address, "")
	fs.StringVar(&sd.Storage, "storage", sd.Storage, "")
	fs.StringVar(&sd.Type, "type", sd.Type, "")
	fs.StringVar(&sd.Tier, "tier", sd.Tier, "")
	fs.StringVar(&sd.Title, "title", sd.Title, "")
	fs.IntVar(&sd.Size, "size", sd.Size, "")
	err = fs.Parse(args)
	if err != nil {
		return nil, err
	}

	if sd.Action != request.CreateServerStorageDeviceActionCreate {
		if sd.Storage == "" {
			return nil, fmt.Errorf("storage UUID or Title must be provided for %s operation", sd.Action)
		}
		strg, err := storage.SearchSingleStorage(sd.Storage, storageSvc)
		if err != nil {
			return nil, err
		}
		sd.Storage = strg.UUID
	}

	if sd.Action == request.CreateServerStorageDeviceActionClone && sd.Title == "" {
		sd.Title = fmt.Sprintf("%s-%s-clone", ui.TruncateText(s.Hostname, 64-7-len(sd.Storage)), sd.Storage)
	}

	if sd.Action == request.CreateServerStorageDeviceActionCreate && sd.Title == "" {
		return nil, fmt.Errorf("title of new storage must be provided")
	}

	return sd, nil
}

func (s *createParams) handleNetwork(in string) (*request.CreateServerInterface, error) {
	var (
		serverInterface             = &request.CreateServerInterface{}
		ipFamily                    string
		ipAddress                   string
		bootable, sourceIPFiltering config.OptionalBoolean
	)
	fs := &pflag.FlagSet{}
	args, err := commands.Parse(in)
	if err != nil {
		return nil, err
	}
	fs.StringVar(&ipFamily, "family", ipFamily, "")
	fs.StringVar(&serverInterface.Type, "type", serverInterface.Type, "")
	fs.StringVar(&serverInterface.Network, "network", "", "")
	fs.StringVar(&ipAddress, "ip-address", "", "")
	config.AddEnableDisableFlags(fs, &bootable, "bootable", "")
	config.AddEnableDisableFlags(fs, &sourceIPFiltering, "source-ip-filtering", "")
	err = fs.Parse(args)
	if err != nil {
		return nil, err
	}
	if serverInterface.Type == "" {
		return nil, fmt.Errorf("network type is required")
	}
	if ipAddress != "" {
		parsedFamily, err := ipaddress.GetFamily(ipAddress)
		if err != nil {
			return nil, err
		}
		if ipFamily != "" && ipFamily != parsedFamily {
			return nil, fmt.Errorf("ip family mismatch: %v != %v", ipFamily, parsedFamily)
		}
		ipFamily = parsedFamily
	} else if ipFamily == "" {
		ipFamily = defaultIPAddressFamily
	}

	serverInterface.Bootable = bootable.AsUpcloudBoolean()
	serverInterface.SourceIPFiltering = sourceIPFiltering.AsUpcloudBoolean()
	serverInterface.IPAddresses = append(serverInterface.IPAddresses,
		request.CreateServerIPAddress{
			Family:  ipFamily,
			Address: ipAddress,
		})

	return serverInterface, nil
}

func (s *createParams) handleSSHKey() error {
	var allSSHKeys []string
	for _, keyOrFile := range s.sshKeys {
		if strings.HasPrefix(keyOrFile, "ssh-") {
			if _, _, _, _, err := ssh.ParseAuthorizedKey([]byte(keyOrFile)); err != nil {
				return fmt.Errorf("invalid ssh key %q: %v", keyOrFile, err)
			}
			allSSHKeys = append(allSSHKeys, keyOrFile)
			continue
		}
		f, err := os.Open(keyOrFile)
		if err != nil {
			return err
		}
		rdr := bufio.NewScanner(f)
		for rdr.Scan() {
			if _, _, _, _, err := ssh.ParseAuthorizedKey(rdr.Bytes()); err != nil {
				_ = f.Close()
				return fmt.Errorf("invalid ssh key %q in file %s: %v", rdr.Text(), keyOrFile, err)
			}
			allSSHKeys = append(allSSHKeys, rdr.Text())
		}
		_ = f.Close()
	}
	s.LoginUser.SSHKeys = allSSHKeys
	return nil
}

type createCommand struct {
	*commands.BaseCommand
	params         createParams
	firewall       config.OptionalBoolean
	metadata       config.OptionalBoolean
	remoteAccess   config.OptionalBoolean
	createPassword config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	s.params = createParams{CreateServerRequest: request.CreateServerRequest{}}
	def := defaultCreateParams
	fs.IntVar(&s.params.AvoidHost, "avoid-host", def.AvoidHost, "Use this to make sure VMs do not reside on specific host. Refers to value from host -attribute. Useful when building HA-environments.")
	fs.IntVar(&s.params.Host, "host", def.Host, "Use this to start a VM on a specific private cloud host. Refers to value from host -attribute. Only available in private clouds.")
	fs.StringVar(&s.params.BootOrder, "boot-order", def.BootOrder, "The boot device order, disk / cdrom / network or comma separated combination.")
	fs.StringVar(&s.params.UserData, "user-data", def.UserData, "Defines URL for a server setup script, or the script body itself.")
	fs.IntVar(&s.params.CoreNumber, "cores", def.CoreNumber, "Number of cores. Use only when defining a flexible (\"custom\") plan.")
	fs.IntVar(&s.params.MemoryAmount, "memory", def.MemoryAmount, "Memory amount in MiB. Use only when defining a flexible (\"custom\") plan.")
	fs.StringVar(&s.params.Title, "title", def.Title, "A short, informational description.")
	fs.StringVar(&s.params.Hostname, "hostname", def.Hostname, "Server hostname.")
	fs.StringVar(&s.params.Plan, "plan", def.Plan, "Server plan name. See \"server plans\" command for valid plans. Set --cores and --memory for a flexible plan.")
	fs.StringVar(&s.params.os, "os", def.os, "Server OS to use (will be the first storage device). Set to empty to fully customise the storages.")
	fs.IntVar(&s.params.osStorageSize, "os-storage-size", def.osStorageSize, "OS storage size in GiB. This is only applicable if `os` is also set. Zero value makes the disk equal to the minimum size of the template.")
	fs.StringVar(&s.params.Zone, "zone", def.Zone, "Zone where to create the server.")
	fs.StringVar(&s.params.PasswordDelivery, "password-delivery", def.PasswordDelivery, "Defines how password is delivered. Available: email, sms")
	fs.StringVar(&s.params.SimpleBackup, "simple-backup", def.SimpleBackup, "Simple backup rule. Format (HHMM,{dailies,weeklies,monthlies}). Example: 2300,dailies")
	fs.StringVar(&s.params.TimeZone, "time-zone", def.TimeZone, "Time zone to set the RTC to.")
	fs.StringVar(&s.params.VideoModel, "video-model", def.VideoModel, "Video interface model of the server. Available: vga, cirrus")
	config.AddEnableOrDisableFlag(fs, &s.firewall, def.firewall, "firewall", "firewall")
	// ----------- fs.BoolVar(&s.params.firewall, "firewall", def.firewall, "Enables the firewall. You can manage firewall rules with the firewall command.")
	config.AddEnableOrDisableFlag(fs, &s.metadata, def.metadata, "metadata", "metadata service")
	// ----------- fs.BoolVar(&s.params.metadata, "metadata", def.metadata, "Enable metadata service.")
	fs.StringArrayVar(&s.params.storages, "storage", def.storages, "A storage connected to the server, multiple can be declared.\nUsage: --storage action=attach,storage=01000000-0000-4000-8000-000020010301,type=cdrom")
	fs.StringArrayVar(&s.params.networks, "network", def.networks, "A network interface for the server, multiple can be declared.\nUsage: --network family=IPv4,type=public\n\n--network type=private,network=037a530b-533e-4cef-b6ad-6af8094bb2bc,ip-address=10.0.0.1")
	config.AddToggleFlag(fs, &s.createPassword, "create-password", def.createPassword, "Create an admin password.")
	fs.StringVar(&s.params.username, "username", def.username, "Admin account username.")
	fs.StringSliceVar(&s.params.sshKeys, "ssh-keys", def.sshKeys, "Add one or more SSH keys to the admin account. Accepted values are SSH public keys or filenames from where to read the keys.")
	config.AddEnableOrDisableFlag(fs, &s.remoteAccess, def.remoteAccess, "remote-access", "remote access")
	// ----------- fs.BoolVar(&s.params.remoteAccess, "remote-access-enabled", def.remoteAccess, "Enables or disables the remote access.")
	fs.StringVar(&s.params.RemoteAccessType, "remote-access-type", def.RemoteAccessType, "Set a remote access type. Available: vnc, spice")
	fs.StringVar(&s.params.RemoteAccessPassword, "remote-access-password", def.RemoteAccessPassword, "Defines the remote access password.")
	s.AddFlags(fs)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	fmt.Printf("\033[31m Plan: %s | Cores: %d | Mem: %d\033[0m", s.params.Plan, s.params.CoreNumber, s.params.MemoryAmount)

	if s.params.Hostname == "" || s.params.Zone == "" {
		return nil, fmt.Errorf("hostname, zone and some password delivery method are required")
	}
	if s.params.os == defaultCreateParams.os && s.params.PasswordDelivery == "none" && s.params.sshKeys == nil {
		return nil, fmt.Errorf("a password-delivery method, ssh-keys or a custom image must be specified")
	}

	if !s.createPassword.Value() && s.params.PasswordDelivery != "none" {
		_ = s.createPassword.Set("true")
	}

	if s.params.Title == "" {
		s.params.Title = s.params.Hostname
	}

	if s.params.CoreNumber != 0 || s.params.MemoryAmount != 0 || s.params.Plan == customPlan {
		if s.params.CoreNumber == 0 || s.params.MemoryAmount == 0 {
			return nil, fmt.Errorf("both --cores and --memory must be defined for custom plans")
		}

		if s.params.Plan != "" && s.params.Plan != customPlan {
			return nil, fmt.Errorf("--plan needs to be 'custom' when --cores and --memory are specified")
		}
	}

	serverSvc := exec.Server()
	planSvc := exec.Plan()
	storageSvc := exec.Storage()
	msg := fmt.Sprintf("creating server %v", s.params.Hostname)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()

	if err := s.params.processParams(planSvc, storageSvc); err != nil {
		return nil, err
	}

	req := s.params.CreateServerRequest
	// TODO: refactor when go-api parameter is refactored
	if s.firewall.Value() {
		req.Firewall = "on"
	}
	req.Metadata = s.metadata.AsUpcloudBoolean()
	req.RemoteAccessEnabled = s.remoteAccess.AsUpcloudBoolean()
	if s.createPassword.Value() {
		req.LoginUser.CreatePassword = "yes"
	}

	logline.SetMessage(fmt.Sprintf("%s: creating network interfaces", msg))
	var iFaces []request.CreateServerInterface
	for _, network := range s.params.networks {
		_interface, err := s.params.handleNetwork(network)
		if err != nil {
			return nil, err
		}
		iFaces = append(iFaces, *_interface)
	}

	logline.SetMessage(fmt.Sprintf("%s: creating storage devices", msg))
	for _, strg := range s.params.storages {
		strg, err := s.params.handleStorage(strg, storageSvc)
		if err != nil {
			return nil, err
		}
		req.StorageDevices = append(req.StorageDevices, *strg)
	}

	if err := s.params.handleSSHKey(); err != nil {
		return nil, err
	}

	if len(iFaces) > 0 {
		req.Networking = &request.CreateServerNetworking{Interfaces: iFaces}
	}

	res, err := serverSvc.CreateServer(&req)
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: request sent", msg))
	logline.MarkDone()

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
		{Title: "IP Addresses", Value: res.IPAddresses, Format: formatIPAddresses},
	}}, nil
}

func formatIPAddresses(val interface{}) (text.Colors, string, error) {
	if ipAddresses, ok := val.(upcloud.IPAddressSlice); ok {
		strs := make([]string, len(ipAddresses))
		for i, ipa := range ipAddresses {
			strs[i] = ipa.Address
		}
		return nil, strings.Join(strs, ",\n"), nil
	}
	return nil, fmt.Sprint(val), nil
}
