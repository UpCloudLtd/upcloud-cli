package storage

import (
	"fmt"
	"sync/atomic"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/cobra"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a storage"),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	service *service.Service
}

func (s *deleteCommand) initService() {
	if s.service == nil {
		s.service = upapi.Service(s.Config())
	}
}

func (s *deleteCommand) InitCommand() {
	s.ArgCompletion(func(toComplete string) ([]string, cobra.ShellCompDirective) {
		s.initService()
		storages, err := s.service.GetStorages(&request.GetStoragesRequest{Access: upcloud.StorageAccessPrivate})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var vals []string
		for _, v := range storages.Storages {
			vals = append(vals, v.UUID, v.Title)
		}
		return commands.MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
	})
	s.SetPositionalArgHelp("<uuidOrTitle ...>")
}

func (s *deleteCommand) MakeExecuteCommand() func(args []string) error {
	return func(args []string) error {
		s.initService()
		if len(args) < 1 {
			return fmt.Errorf("server hostname, title or uuid is required")
		}
		var (
			deleteStorages []*upcloud.Storage
		)
		for _, v := range args {
			storage, err := searchStorage(&cachedStorages, s.service, v, false)
			if err != nil {
				return err
			}
			deleteStorages = append(deleteStorages, storage)
		}
		var numOk int64
		handler := func(idx int, e *ui.LogEntry) {
			storage := deleteStorages[idx]
			msg := fmt.Sprintf("Deleting %q (%s)", storage.Title, storage.UUID)
			e.SetMessage(msg)
			e.Start()
			var err error
			err = s.service.DeleteStorage(&request.DeleteStorageRequest{UUID: storage.UUID})
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				atomic.AddInt64(&numOk, 1)
				e.SetMessage(fmt.Sprintf("%s: done", msg))
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(deleteStorages),
			MaxConcurrentTasks: maxStorageActions,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)

		if int(numOk) < len(deleteStorages) {
			return fmt.Errorf("number of storages failed to delete: %d", len(deleteStorages)-int(numOk))
		}
		return s.HandleOutput(deleteStorages)
	}
}

func (s *deleteCommand) HandleOutput(out interface{}) error {
	results := out.([]*upcloud.Storage)
	var uuids []string
	for _, res := range results {
		uuids = append(uuids, res.UUID)
	}

	if !s.Config().OutputHuman() {
		return s.BaseCommand.HandleOutput(uuids)
	}
	return nil
}
