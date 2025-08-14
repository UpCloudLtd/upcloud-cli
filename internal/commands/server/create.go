package server

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/ipaddress"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/storage"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const defaultIPAddressFamily = upcloud.IPAddressFamilyIPv4

// CreateCommand creates the "server create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a new server",
			"upctl server create --title myapp --zone fi-hel1 --hostname myapp --ssh-keys ~/.ssh/id_*.pub",
			"upctl server create --wait --title myapp --zone fi-hel1 --hostname myapp --ssh-keys ~/.ssh/id_*.pub",
			"upctl server create --title \"My Server\" --zone fi-hel1 --hostname myapp --ssh-keys ~/.ssh/id_*.pub",
			"upctl server create --zone fi-hel1 --hostname myapp --ssh-keys ~/.ssh/id_*.pub --plan 2xCPU-4GB",
			"upctl server create --zone fi-hel1 --hostname myapp --ssh-keys ~/.ssh/id_*.pub --plan custom --cores 2 --memory 4096",
			"upctl server create --zone fi-hel1 --hostname myapp --password-delivery email --os \"Debian GNU/Linux 10 (Buster)\" --server-group a4643646-8342-4324-4134-364138712378",
			"upctl server create --zone fi-hel1 --hostname myapp --ssh-keys ~/.ssh/id_*.pub --network type=private,network=037a530b-533e-4cef-b6ad-6af8094bb2bc,ip-address=10.0.0.1",
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
	os:             "Ubuntu Server 24.04 LTS (Noble Numbat)",
	osStorageSize:  0,
	sshKeys:        nil,
	username:       "",
	createPassword: false,
}

type createParams struct {
	request.CreateServerRequest
	firewall           bool
	metadata           bool
	os                 string
	osStorageSize      int
	osStorageEncrypted config.OptionalBoolean

	labels   []string
	storages []string
	networks []string

	sshKeys        []string
	username       string
	createPassword bool
	remoteAccess   bool
}

func (s *createParams) processParams(exec commands.Executor) error {
	if s.os != "" {
		var osStorage *upcloud.Storage

		osStorage, err := storage.SearchSingleStorage(s.os, exec)
		if err != nil {
			return err
		}

		plans, err := exec.All().GetPlans(exec.Context())
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

		// Enable metadata service for cloud-init templates. Leave empty for other templates to use value defined by user.
		if osStorage.TemplateType == upcloud.StorageTemplateTypeCloudInit {
			s.Metadata = upcloud.True
		}

		s.StorageDevices = append(s.StorageDevices, request.CreateServerStorageDevice{
			Action:  "clone",
			Address: "virtio",
			Storage: osStorage.UUID,
			Title:   fmt.Sprintf("%s-OS", ui.TruncateText(s.Hostname, 64-7)),
			Size:    size,
			Type:    upcloud.StorageTypeDisk,
		})
	}

	if s.osStorageSize != 0 {
		s.StorageDevices[0].Size = s.osStorageSize
	}

	if s.osStorageEncrypted.Value() {
		s.StorageDevices[0].Encrypted = s.osStorageEncrypted.AsUpcloudBoolean()
	}

	if s.LoginUser == nil {
		s.LoginUser = &request.LoginUser{}
	}
	s.LoginUser.CreatePassword = "no"
	if s.username != "" {
		s.LoginUser.Username = s.username
	}

	if len(s.labels) > 0 {
		labelSlice, err := labels.StringsToUpCloudLabelSlice(s.labels)
		if err != nil {
			return err
		}

		s.Labels = labelSlice
	}
	return nil
}

func (s *createParams) handleStorage(in string, exec commands.Executor) (*request.CreateServerStorageDevice, error) {
	var encryptedRaw string

	sd := &request.CreateServerStorageDevice{}
	fs := &pflag.FlagSet{}
	args, err := commands.Parse(in)
	if err != nil {
		return nil, err
	}
	fs.StringVar(&sd.Action, "action", sd.Action, "")
	fs.StringVar(&sd.Address, "address", sd.Address, "")
	fs.StringVar(&encryptedRaw, "encrypt", encryptedRaw, "")
	fs.StringVar(&sd.Storage, "storage", sd.Storage, "")
	fs.StringVar(&sd.Type, "type", sd.Type, "")
	fs.StringVar(&sd.Tier, "tier", sd.Tier, "")
	fs.StringVar(&sd.Title, "title", sd.Title, "")
	fs.IntVar(&sd.Size, "size", sd.Size, "")
	err = fs.Parse(args)
	if err != nil {
		return nil, err
	}

	if encrypted, err := commands.BoolFromString(encryptedRaw); err == nil {
		sd.Encrypted = *encrypted
	}

	if sd.Action != request.CreateServerStorageDeviceActionCreate {
		if sd.Storage == "" {
			return nil, fmt.Errorf("storage UUID or Title must be provided for %s operation", sd.Action)
		}
		strg, err := storage.SearchSingleStorage(sd.Storage, exec)
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
	allSSHKeys, err := commands.ParseSSHKeys(s.sshKeys)
	if err != nil {
		return err
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
	wait           config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	passwordDeliveries := []string{request.PasswordDeliveryNone, request.PasswordDeliveryEmail, request.PasswordDeliverySMS}

	s.Cobra().Long = commands.WrapLongDescription(`Create a new server

Note that the default template, Ubuntu Server 24.04 LTS (Noble Numbat), only supports SSH key based authentication. Use ` + "`" + `--ssh-keys` + "`" + ` option to provide the keys when creating a server with the default template. The examples below use public key from the ` + "`" + `~/.ssh` + "`" + ` directory. If you want to use different authentication method, use ` + "`" + `--os` + "`" + ` parameter to specify a different template.`)

	fs := &pflag.FlagSet{}
	s.params = createParams{CreateServerRequest: request.CreateServerRequest{}}
	def := defaultCreateParams
	fs.IntVar(&s.params.AvoidHost, "avoid-host", def.AvoidHost, avoidHostDescription)
	fs.StringVar(&s.params.BootOrder, "boot-order", def.BootOrder, "The boot device order, disk / cdrom / network or comma separated combination.")
	fs.IntVar(&s.params.CoreNumber, "cores", def.CoreNumber, "Number of cores. Only allowed if `plan` option is set to \"custom\".")
	config.AddToggleFlag(fs, &s.createPassword, "create-password", def.createPassword, "Create an admin password.")
	config.AddEnableOrDisableFlag(fs, &s.firewall, def.firewall, "firewall", "firewall")
	config.AddEnableOrDisableFlag(fs, &s.metadata, def.metadata, "metadata", "metadata service. The metadata service will be enabled by default, if the selected OS template uses cloud-init and thus requires metadata service")
	config.AddEnableOrDisableFlag(fs, &s.remoteAccess, def.remoteAccess, "remote-access", "remote access")
	fs.IntVar(&s.params.Host, "host", def.Host, hostDescription)
	fs.StringVar(&s.params.Hostname, "hostname", def.Hostname, "Server hostname.")
	fs.StringArrayVar(&s.params.labels, "label", def.labels, "Labels to describe the server in `key=value` format, multiple can be declared.\nUsage: --label env=dev\n\n--label owner=operations")
	fs.IntVar(&s.params.MemoryAmount, "memory", def.MemoryAmount, "Memory amount in MiB. Only allowed if `plan` option is set to \"custom\".")
	fs.StringArrayVar(&s.params.networks, "network", def.networks, "A network interface for the server, multiple can be declared.\nUsage: --network family=IPv4,type=public\n\n--network type=private,network=037a530b-533e-4cef-b6ad-6af8094bb2bc,ip-address=10.0.0.1")
	fs.StringVar(&s.params.os, "os", def.os, "Server OS to use (will be the first storage device). The value should be title or UUID of an either public or private template. Set to empty to fully customise the storages.")
	fs.IntVar(&s.params.osStorageSize, "os-storage-size", def.osStorageSize, "OS storage size in GiB. This is only applicable if `os` is also set. Zero value makes the disk equal to the minimum size of the template.")
	config.AddToggleFlag(fs, &s.params.osStorageEncrypted, "os-storage-encrypt", false, "Encrypt the OS storage. This is only applicable if `os` is also set.")
	fs.StringVar(&s.params.PasswordDelivery, "password-delivery", def.PasswordDelivery, "Defines how password is delivered. Available: "+strings.Join(passwordDeliveries, ", "))
	fs.StringVar(&s.params.Plan, "plan", def.Plan, "Server plan name. See \"server plans\" command for valid plans. Set to \"custom\" and use `cores` and `memory` options for flexible plan.")
	fs.StringVar(&s.params.RemoteAccessPassword, "remote-access-password", def.RemoteAccessPassword, "Defines the remote access password.")
	fs.StringVar(&s.params.RemoteAccessType, "remote-access-type", def.RemoteAccessType, "Set a remote access type. Available: "+strings.Join(remoteAccessTypes, ", "))
	fs.StringVar(&s.params.ServerGroup, "server-group", def.ServerGroup, "UUID of a server group for the server. To remove the server from the group, see `servergroup modify")
	fs.StringVar(&s.params.SimpleBackup, "simple-backup", def.SimpleBackup, simpleBackupDescription)
	fs.StringSliceVar(&s.params.sshKeys, "ssh-keys", def.sshKeys, "Add one or more SSH keys to the admin account. Accepted values are SSH public keys or filenames from where to read the keys.")
	fs.StringArrayVar(&s.params.storages, "storage", def.storages, "A storage connected to the server, multiple can be declared.\nUsage: --storage action=attach,storage=01000000-0000-4000-8000-000020010301,type=cdrom")
	fs.StringVar(&s.params.TimeZone, "time-zone", def.TimeZone, "Time zone to set the RTC to.")
	fs.StringVar(&s.params.Title, "title", def.Title, "A short, informational description.")
	fs.StringVar(&s.params.UserData, "user-data", def.UserData, "Defines URL for a server setup script, or the script body itself.")
	fs.StringVar(&s.params.username, "username", def.username, "Admin account username.")
	fs.StringVar(&s.params.VideoModel, "video-model", def.VideoModel, "Video interface model of the server. Available: "+strings.Join(videoModels, ", "))
	config.AddToggleFlag(fs, &s.wait, "wait", false, "Wait for server to be in started state before returning.")
	fs.StringVar(&s.params.Zone, "zone", def.Zone, namedargs.ZoneDescription("server"))
	// fs.BoolVar(&s.params.firewall, "firewall", def.firewall, "Enables the firewall. You can manage firewall rules with the firewall command.")
	// fs.BoolVar(&s.params.metadata, "metadata", def.metadata, "Enable metadata service.")
	// fs.BoolVar(&s.params.remoteAccess, "remote-access-enabled", def.remoteAccess, "Enables or disables the remote access.")
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("hostname"))
	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("password-delivery", cobra.FixedCompletions(passwordDeliveries, cobra.ShellCompDirectiveNoFileComp)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("remote-access-type", cobra.FixedCompletions(remoteAccessTypes, cobra.ShellCompDirectiveNoFileComp)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("video-model", cobra.FixedCompletions(videoModels, cobra.ShellCompDirectiveNoFileComp)))
	for _, flag := range []string{
		"boot-order", "cores", "hostname", "label", "memory", "network", "os", "os-storage-size", "remote-access-password",
		"simple-backup", "storage", "title", "user-data", "username",
	} {
		commands.Must(s.Cobra().RegisterFlagCompletionFunc(flag, cobra.NoFileCompletions))
	}
}

func (s *createCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("avoid-host", namedargs.CompletionFunc(completion.HostID{}, cfg)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("host", namedargs.CompletionFunc(completion.HostID{}, cfg)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("plan", namedargs.CompletionFunc(completion.ServerPlan{}, cfg)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("server-group", namedargs.CompletionFunc(completion.ServerGroupUUID{}, cfg)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("time-zone", namedargs.CompletionFunc(completion.TimeZone{}, cfg)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("zone", namedargs.CompletionFunc(completion.Zone{}, cfg)))
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
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

		if s.params.Plan != customPlan {
			return nil, fmt.Errorf("--plan needs to be 'custom' when --cores and --memory are specified")
		}
	}

	svc := exec.All()
	msg := fmt.Sprintf("Creating server %v", s.params.Hostname)
	exec.PushProgressStarted(msg)

	if err := s.params.processParams(exec); err != nil {
		return nil, err
	}

	req := s.params.CreateServerRequest
	// TODO: refactor when go-api parameter is refactored
	if s.firewall.Value() {
		req.Firewall = "on"
	}
	if req.Metadata.Empty() {
		req.Metadata = s.metadata.AsUpcloudBoolean()
	}
	req.RemoteAccessEnabled = s.remoteAccess.AsUpcloudBoolean()
	if s.createPassword.Value() {
		req.LoginUser.CreatePassword = "yes"
	}

	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("%s: creating network interfaces", msg))
	var iFaces []request.CreateServerInterface
	for _, network := range s.params.networks {
		_interface, err := s.params.handleNetwork(network)
		if err != nil {
			return nil, err
		}
		iFaces = append(iFaces, *_interface)
	}

	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("%s: creating storage devices", msg))
	for _, strg := range s.params.storages {
		strg, err := s.params.handleStorage(strg, exec)
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

	exec.PushProgressUpdateMessage(msg, msg)
	res, err := svc.CreateServer(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if s.wait.Value() {
		waitForServerState(res.UUID, upcloud.ServerStateStarted, exec, msg)
	} else {
		exec.PushProgressSuccess(msg)
	}

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
		{Title: "IP Addresses", Value: res, Format: formatCreateIPAddresses},
	}}, nil
}

func formatCreateIPAddresses(val interface{}) (text.Colors, string, error) {
	server, ok := val.(*upcloud.ServerDetails)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse IP addresses from %T, expected *upcloud.ServerDetails", val)
	}

	// Store addresses in map keys to avoid duplicate addresses
	addresses := make(map[string]bool)

	// Get public addresses from ip_addresses list
	// Public and utility interfaces created by default (no --network parameters) are only listed here
	for _, ipa := range server.IPAddresses {
		addresses[ipa.Address] = true
	}

	// Get public and private addresses from networking.interfaces list
	// Public and utility interfaces created with --network are also listed here
	for _, iface := range server.Networking.Interfaces {
		for _, ipa := range iface.IPAddresses {
			addresses[ipa.Address] = true
		}
	}

	var strs []string
	for ipa := range addresses {
		strs = append(strs, ui.DefaultAddressColours.Sprint(ipa))
	}

	return nil, strings.Join(strs, ",\n"), nil
}
