package storage

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"io"
	"math"
	"sync"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/ui"
)

func ShowCommand(serverSvc service.Server, storageSvc service.Storage) commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show storage details"),
		serverSvc:   serverSvc,
		storageSvc:  storageSvc,
	}
}

type showCommand struct {
	*commands.BaseCommand
	serverSvc     service.Server
	storageSvc    service.Storage
	storageImport *upcloud.StorageImportDetails
}

type commandResponseHolder struct {
	storageDetails *upcloud.StorageDetails
	storageImport  *upcloud.StorageImportDetails
	servers        []upcloud.Server
	storages       []upcloud.Storage
}

func (s *showCommand) InitCommand() {
	s.ArgCompletion(func(toComplete string) ([]string, cobra.ShellCompDirective) {
		storages, err := s.storageSvc.GetStorages(&request.GetStoragesRequest{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var vals []string
		for _, v := range storages.Storages {
			vals = append(vals, v.UUID, v.Title)
		}
		return commands.MatchStringPrefix(vals, toComplete, false), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
}

func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		var storages []upcloud.Storage
		var storageImport *upcloud.StorageImportDetails
		var servers []upcloud.Server

		if len(args) < 1 {
			return nil, fmt.Errorf("storage title or uuid is required")
		}
		storage, err := searchStorage(&storages, s.storageSvc, args[0], true)
		if err != nil {
			return nil, err
		}
		var (
			wg                      sync.WaitGroup
			storageImportDetailsErr error
		)
		wg.Add(1)
		go func() {
			defer wg.Done()
			storageImport, storageImportDetailsErr = s.storageSvc.GetStorageImportDetails(
				&request.GetStorageImportDetailsRequest{UUID: storage[0].UUID})

			if ucErr, ok := storageImportDetailsErr.(*upcloud.Error); ok {
				if ucErr.ErrorCode == "STORAGE_IMPORT_NOT_FOUND" {
					storageImportDetailsErr = nil
				}
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			if resp, err := s.serverSvc.GetServers(); err == nil {
				servers = resp.Servers
			}
		}()
		storageDetails, err := s.storageSvc.GetStorageDetails(&request.GetStorageDetailsRequest{UUID: storage[0].UUID})
		if err != nil {
			return nil, err
		}
		wg.Wait()
		if storageImportDetailsErr != nil {
			return nil, storageImportDetailsErr
		}
		return &commandResponseHolder{storageDetails, storageImport, servers, storages}, nil
	}
}

func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
	resp := out.(*commandResponseHolder)
	storage := resp.storageDetails
	storageImport := resp.storageImport

	formatBool := func(v bool) interface{} {
		if v {
			return ui.DefaultBooleanColoursTrue.Sprint("yes")
		}
		return ui.DefaultBooleanColoursFalse.Sprint("no")
	}
	storageByUuid := make(map[string]upcloud.Storage)
	formatStorageReferenceUuid := func(uuid string) string {
		if uuid == "" {
			return ""
		}
		if len(storageByUuid) == 0 {
			for _, v := range resp.storages {
				storageByUuid[v.UUID] = v
			}
		}
		if v, ok := storageByUuid[uuid]; ok {
			return fmt.Sprintf("%s (%s)", v.Title, ui.DefaultUuidColours.Sprint(uuid))
		}
		return ui.DefaultUuidColours.Sprint(uuid)
	}
	rowTransformer := func(row table.Row) table.Row {
		if v, ok := row[len(row)-1].(upcloud.Boolean); ok {
			row[len(row)-1] = formatBool(v.Bool())
		}
		if v, ok := row[len(row)-1].(time.Time); ok {
			row[len(row)-1] = ui.FormatTime(v)
		}
		if v, ok := row[len(row)-1].(float64); ok {
			if _, frac := math.Modf(v); frac != 0 {
				row[len(row)-1] = fmt.Sprintf("%.2f", v)
			}
		}
		return row
	}

	dMain := ui.NewListLayout(ui.ListLayoutDefault)

	// Common details
	{
		dCommon := ui.NewDetailsView()
		dCommon.SetRowTransformer(rowTransformer)
		dCommon.AppendRows([]table.Row{
			{"UUID:", ui.DefaultUuidColours.Sprint(storage.UUID)},
			{"Title:", storage.Title},
			{"Zone:", storage.Zone},
			{"State:", StateColour(storage.State).Sprint(storage.State)},
			{"Size (GiB):", storage.Size},
			{"Type:", storage.Type},
			{"Tier:", storage.Tier},
			{"License:", storage.License},
			{"Created:", storage.Created},
			{"Origin:", formatStorageReferenceUuid(storage.Origin)},
		})
		dMain.AppendSection("Common:", dCommon.Render())
	}

	// Servers
	if len(storage.ServerUUIDs) > 0 {
		tServers := ui.NewDataTable("UUID", "Title", "Hostname", "State")
		serversByUuid := make(map[string]upcloud.Server)
		for _, v := range resp.servers {
			serversByUuid[v.UUID] = v
		}
		for _, uuid := range storage.ServerUUIDs {
			tServers.AppendRow(table.Row{
				ui.DefaultUuidColours.Sprint(uuid),
				serversByUuid[uuid].Title,
				serversByUuid[uuid].Hostname,
				server.StateColour(serversByUuid[uuid].State).Sprint(serversByUuid[uuid].State),
			})
		}
		dMain.AppendSection("Servers:", ui.WrapWithListLayout(tServers.Render(), ui.ListLayoutNestedTable).Render())
	} else {
		dMain.AppendSection("Servers:", "(no servers using this storage)")
	}

	// Backups
	{
		dBackups := ui.NewDetailsView()
		dBackups.SetRowSpacing(true)
		if storage.BackupRule != nil && storage.BackupRule.Interval != "" {
			dBackupRule := ui.NewDetailsView()
			dBackupRule.AppendRows([]table.Row{
				{"Interval:", storage.BackupRule.Interval},
				{"Time:", storage.BackupRule.Time},
				{"Retention:", storage.BackupRule.Retention},
			})
			dMain.AppendSection("Backup Rule:", dBackupRule.Render())
		} else if storage.BackupRule != nil {
			dMain.AppendSection("Backup Rule:", "(no backup rule configured)")
		}

		if len(storage.BackupUUIDs) > 0 {
			if len(storageByUuid) == 0 {
				for _, v := range resp.storages {
					storageByUuid[v.UUID] = v
				}
			}
			tBackups := ui.NewDataTable("UUID", "Title", "Created")
			for _, uuid := range storage.BackupUUIDs {
				tBackups.AppendRow(table.Row{
					ui.DefaultUuidColours.Sprint(uuid),
					storageByUuid[uuid].Title,
					storageByUuid[uuid].Created,
				})
			}
			dMain.AppendSection("Available Backups:", ui.WrapWithListLayout(tBackups.Render(), ui.ListLayoutNestedTable).Render())
		} else {
			dMain.AppendSection("Available Backups:", "(no backups are available)")
		}
	}

	// Storage import
	if storageImport != nil {
		dStorageImport := ui.NewDetailsView()
		dStorageImport.SetRowTransformer(rowTransformer)
		dStorageImport.AppendRows([]table.Row{
			{"State:", ImportStateColour(storageImport.State).Sprint(storageImport.State)},
			{"Source:", storageImport.Source},
		})
		switch storageImport.Source {
		case upcloud.StorageImportSourceHTTPImport:
			dStorageImport.AppendRow(table.Row{"Source Location", storageImport.SourceLocation})
		case upcloud.StorageImportSourceDirectUpload:
			dStorageImport.AppendRow(table.Row{"Upload URL", storageImport.DirectUploadURL})
		}
		dStorageImport.AppendRows([]table.Row{
			{"Content Length:", ui.FormatBytes(storageImport.ClientContentLength)},
			{"Read:", ui.FormatBytes(storageImport.ReadBytes)},
			{"Written:", ui.FormatBytes(storageImport.WrittenBytes)},
			{"SHA256 Checksum:", storageImport.SHA256Sum},
		})
		if storageImport.ErrorCode != "" {
			dStorageImport.AppendRows([]table.Row{
				{"Error:", ui.DefaultErrorColours.Sprintf("%s\n%s",
					storageImport.ErrorCode, storageImport.ErrorMessage)},
			})
		}
		dStorageImport.AppendRows([]table.Row{
			{"Content Type:", storageImport.ClientContentType},
			{"Created:", ui.FormatTime(storageImport.Created)},
			{"Completed:", ui.FormatTime(storageImport.Completed)},
		})
		dMain.AppendSection("Import:", dStorageImport.Render())
	}

	_, _ = fmt.Fprintln(writer, dMain.Render())
	return nil
}
