package storage

import (
	"encoding/json"
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"io"
	"math"
	"sync"
	"time"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/table"
)

// ShowCommand creates the "storage show" command
func ShowCommand(serverSvc service.Server, storageSvc service.Storage) commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show storage details"),
		serverSvc:   serverSvc,
		storageSvc:  storageSvc,
	}
}

type showCommand struct {
	*commands.BaseCommand
	serverSvc  service.Server
	storageSvc service.Storage
}

type commandResponseHolder struct {
	storageDetails *upcloud.StorageDetails
	storageImport  *upcloud.StorageImportDetails
	servers        []upcloud.Server
	storages       []upcloud.Storage
}

// MarshalJSON implements json.Marshaler
func (c *commandResponseHolder) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.storageDetails)
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getStorageArgumentCompletionFunction(s.storageSvc))
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		var storages []upcloud.Storage
		var storageImport *upcloud.StorageImportDetails
		var servers []upcloud.Server

		if len(args) != 1 {
			return nil, fmt.Errorf("one storage title or uuid is required")
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

// HandleOutput implements Command.HandleOutput
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
	storageByUUID := make(map[string]upcloud.Storage)
	formatStorageReferenceUUID := func(uuid string) string {
		if uuid == "" {
			return ""
		}
		if len(storageByUUID) == 0 {
			for _, v := range resp.storages {
				storageByUUID[v.UUID] = v
			}
		}
		if v, ok := storageByUUID[uuid]; ok {
			return fmt.Sprintf("%s (%s)", v.Title, ui.DefaultUUUIDColours.Sprint(uuid))
		}
		return ui.DefaultUUUIDColours.Sprint(uuid)
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
		dCommon.Append(
			table.Row{"UUID:", ui.DefaultUUUIDColours.Sprint(storage.UUID)},
			table.Row{"Title:", storage.Title},
			table.Row{"Zone:", storage.Zone},
			table.Row{"State:", storageStateColor(storage.State).Sprint(storage.State)},
			table.Row{"Size (GiB):", storage.Size},
			table.Row{"Type:", storage.Type},
			table.Row{"Tier:", storage.Tier},
			table.Row{"Licence:", storage.License},
			table.Row{"Created:", storage.Created},
			table.Row{"Origin:", formatStorageReferenceUUID(storage.Origin)},
		)
		dMain.AppendSection("Common:", dCommon.Render())
	}

	// Servers
	if len(storage.ServerUUIDs) > 0 {
		tServers := ui.NewDataTable("UUID", "Title", "Hostname", "State")
		serversByUUID := make(map[string]upcloud.Server)
		for _, v := range resp.servers {
			serversByUUID[v.UUID] = v
		}
		for _, uuid := range storage.ServerUUIDs {
			tServers.Append(table.Row{
				ui.DefaultUUUIDColours.Sprint(uuid),
				serversByUUID[uuid].Title,
				serversByUUID[uuid].Hostname,
				commands.StateColour(serversByUUID[uuid].State).Sprint(serversByUUID[uuid].State),
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
			dBackupRule.Append(
				table.Row{"Interval:", storage.BackupRule.Interval},
				table.Row{"Time:", storage.BackupRule.Time},
				table.Row{"Retention:", storage.BackupRule.Retention},
			)
			dMain.AppendSection("Backup Rule:", dBackupRule.Render())
		} else if storage.BackupRule != nil {
			dMain.AppendSection("Backup Rule:", "(no backup rule configured)")
		}

		if len(storage.BackupUUIDs) > 0 {
			if len(storageByUUID) == 0 {
				for _, v := range resp.storages {
					storageByUUID[v.UUID] = v
				}
			}
			tBackups := ui.NewDataTable("UUID", "Title", "Created")
			for _, uuid := range storage.BackupUUIDs {
				tBackups.Append(table.Row{
					ui.DefaultUUUIDColours.Sprint(uuid),
					storageByUUID[uuid].Title,
					storageByUUID[uuid].Created,
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
		dStorageImport.Append(
			table.Row{"State:", importStateColor(storageImport.State).Sprint(storageImport.State)},
			table.Row{"Source:", storageImport.Source},
		)
		switch storageImport.Source {
		case upcloud.StorageImportSourceHTTPImport:
			dStorageImport.Append(table.Row{"Source Location", storageImport.SourceLocation})
		case upcloud.StorageImportSourceDirectUpload:
			dStorageImport.Append(table.Row{"Upload URL", storageImport.DirectUploadURL})
		}
		dStorageImport.Append(
			table.Row{"Content Length:", ui.FormatBytes(storageImport.ClientContentLength)},
			table.Row{"Read:", ui.FormatBytes(storageImport.ReadBytes)},
			table.Row{"Written:", ui.FormatBytes(storageImport.WrittenBytes)},
			table.Row{"SHA256 Checksum:", storageImport.SHA256Sum},
		)
		if storageImport.ErrorCode != "" {
			dStorageImport.Append(
				table.Row{"Error:", ui.DefaultErrorColours.Sprintf("%s\n%s",
					storageImport.ErrorCode, storageImport.ErrorMessage)},
			)
		}
		dStorageImport.Append(
			table.Row{"Content Type:", storageImport.ClientContentType},
			table.Row{"Created:", ui.FormatTime(storageImport.Created)},
			table.Row{"Completed:", ui.FormatTime(storageImport.Completed)},
		)
		dMain.AppendSection("Import:", dStorageImport.Render())
	}

	_, _ = fmt.Fprintln(writer, dMain.Render())
	return nil
}
