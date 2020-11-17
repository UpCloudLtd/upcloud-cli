package storage

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/interfaces"
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

func ShowCommand(service interfaces.StorageServer) commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show storage details"),
		service:     service,
	}
}

type showCommand struct {
	*commands.BaseCommand
	service       interfaces.StorageServer
	storageImport *upcloud.StorageImportDetails
}

func (s *showCommand) InitCommand() {
	s.ArgCompletion(func(toComplete string) ([]string, cobra.ShellCompDirective) {
		storages, err := s.service.GetStorages(&request.GetStoragesRequest{})
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
		if len(args) < 1 {
			return nil, fmt.Errorf("storage title or uuid is required")
		}
		storage, err := searchStorage(&cachedStorages, s.service, args[0], true)
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
			s.storageImport, storageImportDetailsErr = s.service.GetStorageImportDetails(
				&request.GetStorageImportDetailsRequest{UUID: storage.UUID})
			if ucErr, ok := storageImportDetailsErr.(*upcloud.Error); ok {
				if ucErr.ErrorCode == "STORAGE_IMPORT_NOT_FOUND" {
					storageImportDetailsErr = nil
				}
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			if servers, err := s.service.GetServers(); err == nil {
				cachedServers = servers.Servers
			}
		}()
		storageDetails, err := s.service.GetStorageDetails(&request.GetStorageDetailsRequest{UUID: storage.UUID})
		if err != nil {
			return nil, err
		}
		wg.Wait()
		if storageImportDetailsErr != nil {
			return nil, storageImportDetailsErr
		}
		return storageDetails, nil
	}
}

func (s *showCommand) HandleOutput(out interface{}) (string, error) {
	storage := out.(*upcloud.StorageDetails)

	dMain := ui.NewDetailsView()
	dMain.SetRowSeparators(true)
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
			for _, v := range cachedStorages {
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

	// Common details
	{
		dCommon := ui.NewDetailsView()
		dCommon.SetRowTransformer(rowTransformer)
		dCommon.AppendRows([]table.Row{
			{"UUID", ui.DefaultUuidColours.Sprint(storage.UUID)},
			{"Title", storage.Title},
			{"Zone", storage.Zone},
			{"State", StateColour(storage.State).Sprint(storage.State)},
			{"Size (GiB)", storage.Size},
			{"Type", storage.Type},
			{"Tier", storage.Tier},
			{"License", storage.License},
			{"Created", storage.Created},
			{"Origin", formatStorageReferenceUuid(storage.Origin)},
		})
		// fmt.Println(dCommon.Render())
		dMain.AppendRow(table.Row{"Common", dCommon.Render()})
	}

	// Servers
	if len(storage.ServerUUIDs) > 0 {
		tServers := ui.NewDataTable("UUID", "Title", "Hostname", "State")
		serversByUuid := make(map[string]upcloud.Server)
		for _, v := range cachedServers {
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
		dMain.AppendRow(table.Row{"Servers", tServers.Render()})
	} else {
		dMain.AppendRow(table.Row{"Servers", "no servers using this storage"})
	}

	// Backups
	{
		dBackups := ui.NewDetailsView()
		dBackups.SetRowSpacing(true)
		if storage.BackupRule != nil && storage.BackupRule.Interval != "" {
			dBackupRule := ui.NewDetailsView()
			dBackupRule.AppendRows([]table.Row{
				{"Interval", storage.BackupRule.Interval},
				{"Time", storage.BackupRule.Time},
				{"Retention", storage.BackupRule.Retention},
			})
			dBackups.AppendRow(table.Row{"Backup Rule", dBackupRule.Render()})
		} else if storage.BackupRule != nil {
			dBackups.AppendRow(table.Row{"Backup Rule", "no backup rule configured"})
		}
		if len(storage.BackupUUIDs) > 0 {
			if len(storageByUuid) == 0 {
				for _, v := range cachedStorages {
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
			dBackups.AppendRow(table.Row{"Available Backups", tBackups.Render()})
		}
		dMain.AppendRow(table.Row{"Backup", dBackups.Render()})
	}

	// Storage import
	if s.storageImport != nil {
		dStorageImport := ui.NewDetailsView()
		dStorageImport.SetRowTransformer(rowTransformer)
		dStorageImport.AppendRows([]table.Row{
			{"State", ImportStateColour(s.storageImport.State).Sprint(s.storageImport.State)},
			{"Source", s.storageImport.Source},
		})
		switch s.storageImport.Source {
		case upcloud.StorageImportSourceHTTPImport:
			dStorageImport.AppendRow(table.Row{"Source Location", s.storageImport.SourceLocation})
		case upcloud.StorageImportSourceDirectUpload:
			dStorageImport.AppendRow(table.Row{"Upload URL", s.storageImport.DirectUploadURL})
		}
		dStorageImport.AppendRows([]table.Row{
			{"Content Length", ui.FormatBytes(s.storageImport.ClientContentLength)},
			{"Read", ui.FormatBytes(s.storageImport.ReadBytes)},
			{"Written", ui.FormatBytes(s.storageImport.WrittenBytes)},
			{"SHA256 Checksum", s.storageImport.SHA256Sum},
		})
		if s.storageImport.ErrorCode != "" {
			dStorageImport.AppendRows([]table.Row{
				{"Error", ui.DefaultErrorColours.Sprintf("%s\n%s",
					s.storageImport.ErrorCode, s.storageImport.ErrorMessage)},
			})
		}
		dStorageImport.AppendRows([]table.Row{
			{"Content Type", s.storageImport.ClientContentType},
			{"Created", ui.FormatTime(s.storageImport.Created)},
			{"Completed", ui.FormatTime(s.storageImport.Completed)},
		})
		dMain.AppendRow(table.Row{"Import", dStorageImport.Render()})
	}

	return dMain.Render(), nil
}
