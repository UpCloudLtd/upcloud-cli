package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

var (
	cachedTemplates []upcloud.Storage
)

func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a server"),
	}
}

var DefaultCreateParams = &createParams{
	CreateServerRequest: request.CreateServerRequest{
		VideoModel: "vga",
		TimeZone:   "UTC",
		Plan:       "1xCPU-2GB",
	},
	firewall:       false,
	metadata:       false,
	os:             "Debian GNU/Linux 10 (Buster)",
	osStorageSize:  0,
	sshKeys:        nil,
	username:       "",
	createPassword: true,
}

type createParams struct {
	request.CreateServerRequest
	firewall      bool
	metadata      bool
	os            string
	osStorageSize int

	sshKeys        []string
	username       string
	createPassword bool
}

func (s *createParams) processParams(srv *service.Service) error {
	if s.os != "" {
		var osStorage *upcloud.Storage
		if len(cachedTemplates) == 0 {
			tpls, err := srv.GetStorages(&request.GetStoragesRequest{
				Type: "template",
			})
			if err != nil {
				return err
			}
			cachedTemplates = tpls.Storages
		}
		for _, tpl := range cachedTemplates {
			if tpl.Title == s.os {
				osStorage = &tpl
				break
			}
			if tpl.UUID == s.os {
				osStorage = &tpl
				break
			}
		}
		size := minStorageSize
		if s.osStorageSize > size {
			size = s.osStorageSize
		}
		if osStorage == nil {
			return fmt.Errorf("no OS storage found with title or uuid %q", s.os)
		}
		s.StorageDevices = append(s.StorageDevices, request.CreateServerStorageDevice{
			Action:  "clone",
			Storage: osStorage.UUID,
			Title:   fmt.Sprintf("%s-osDisk", ui.TruncateText(s.Hostname, 64-7)),
			Size:    size,
			Tier:    upcloud.StorageTierMaxIOPS,
			Type:    upcloud.StorageTypeDisk,
		})
	}
	if s.osStorageSize != 0 {
		s.StorageDevices[0].Size = s.osStorageSize
	}

	if s.firewall {
		s.Firewall = "on"
	}
	if s.metadata {
		s.Metadata = 1
	}
	if s.LoginUser == nil {
		s.LoginUser = &request.LoginUser{}
	}
	s.LoginUser.CreatePassword = "no"
	if s.createPassword {
		s.LoginUser.CreatePassword = "yes"
	}
	if s.username != "" {
		s.LoginUser.Username = s.username
	}
	var allSshKeys []string
	for _, keyOrFile := range s.sshKeys {
		if strings.HasPrefix(keyOrFile, "ssh-") {
			if _, _, _, _, err := ssh.ParseAuthorizedKey([]byte(keyOrFile)); err != nil {
				return fmt.Errorf("invalid ssh key %q: %v", keyOrFile, err)
			}
			allSshKeys = append(allSshKeys, keyOrFile)
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
			allSshKeys = append(allSshKeys, rdr.Text())
		}
		_ = f.Close()
	}
	s.LoginUser.SSHKeys = allSshKeys
	return nil
}

type createCommand struct {
	*commands.BaseCommand
	service           *service.Service
	avoidHost         int
	host              int
	firstCreateServer createParams
	flagSet           *pflag.FlagSet
}

func (s *createCommand) initService() {
	if s.service == nil {
		s.service = upapi.Service(s.Config())
	}
}

func (s *createCommand) createFlags(fs *pflag.FlagSet, dst, def *createParams) {
	fs.IntVar(&dst.CoreNumber, "cores", def.CoreNumber, "Number of cores")
	fs.IntVar(&dst.MemoryAmount, "memory", def.MemoryAmount, "Memory amount in MiB")
	fs.StringVar(&dst.Title, "title", def.Title, "Visible name")
	fs.StringVar(&dst.Hostname, "hostname", def.Hostname, "Hostname")
	fs.StringVar(&dst.Plan, "plan", def.Plan, "Server plan to use. "+
		"Set this to custom to use custom core/memory amounts.")
	fs.StringVar(&dst.os, "os", def.os,
		"Server OS to use (will be the first storage device). Set to empty to fully customise the storages.")
	fs.IntVar(&dst.osStorageSize, "os-storage-size", def.osStorageSize,
		"OS storage size in GiB. This is only applicable if `os` is also set. "+
			"Zero value makes the disk equal to the minimum size of the template.")
	fs.StringVar(&dst.Zone, "zone", def.Zone, "Zone where to create the server")
	fs.StringVar(&dst.PasswordDelivery, "password-delivery", def.PasswordDelivery,
		"If password login is enable set a way how password is delivered.\nAvailable: email,sms")
	fs.StringVar(&dst.SimpleBackup, "simple-backup", def.SimpleBackup,
		"Simple backup rule. Format (HHMM,{dailies,weeklies,monthlies}).\nExample: 2300,dailies")
	fs.StringVar(&dst.TimeZone, "time-zone", def.TimeZone, "Time zone to set the RTC to")
	fs.StringVar(&dst.VideoModel, "video-model", def.VideoModel,
		"Video interface model of the server.\nAvailable: vga,cirrus")
	fs.BoolVar(&dst.firewall, "firewall", def.firewall,
		"Sets the firewall on. You can manage firewall rules with the firewall command")
	fs.BoolVar(&dst.metadata, "metadata", def.metadata, "Enable metadata service")
	fs.BoolVar(&dst.createPassword, "create-password", def.createPassword, "Create a admin password")
	fs.StringVar(&dst.username, "username", def.username, "Admin account username")
	fs.StringSliceVar(&dst.sshKeys, "ssh-keys", def.sshKeys,
		"Add one or more SSH keys to the admin account. Accepted values are SSH public keys or "+
			"filenames from where to read the keys.")

}

func (s *createCommand) InitCommand() {
	s.flagSet = &pflag.FlagSet{}
	s.createFlags(s.flagSet, &s.firstCreateServer, DefaultCreateParams)
	s.AddFlags(s.flagSet)
}

func (s *createCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		s.initService()
		var createServers []request.CreateServerRequest
		if err := s.firstCreateServer.processParams(s.service); err != nil {
			return nil, err
		}
		createServers = append(createServers, s.firstCreateServer.CreateServerRequest)

		// Process additional server create args
		var additionalCreateArgs = make([]string, 0, len(args))
		for i, arg := range args {
			if arg == "--" || i == len(args)-1 {
				if i == len(args)-1 && arg != "--" {
					additionalCreateArgs = append(additionalCreateArgs, arg)
				}
				if len(additionalCreateArgs) > 0 {
					fs := &pflag.FlagSet{}
					dst := createParams{}
					s.createFlags(fs, &dst, &s.firstCreateServer)
					if err := fs.Parse(additionalCreateArgs); err != nil {
						return nil, err
					}
					if err := dst.processParams(s.service); err != nil {
						return nil, err
					}
					createServers = append(createServers, dst.CreateServerRequest)
				}
				additionalCreateArgs = additionalCreateArgs[:0]
				continue
			}
			additionalCreateArgs = append(additionalCreateArgs, arg)
		}

		var (
			mu             sync.Mutex
			numOk          int
			createdServers []*upcloud.ServerDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			server := createServers[idx]
			msg := fmt.Sprintf("Creating server %q", server.Hostname)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.CreateServer(&server)
			if err == nil {
				e.SetMessage(fmt.Sprintf("%s: server starting", msg))
				details, err = WaitForServerState(s.service, details.UUID, upcloud.ServerStateStarted,
					s.Config().ClientTimeout())
			}
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				e.SetMessage(fmt.Sprintf("%s: done", msg))
				mu.Lock()
				numOk++
				createdServers = append(createdServers, details)
				mu.Unlock()
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(createServers),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(createServers) {
			return nil, fmt.Errorf("number of servers that failed: %d", len(createServers)-numOk)
		}
		return createdServers, nil
	}
}
