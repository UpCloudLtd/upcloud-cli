package storage

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
	"io"
	"sync"
)

type showImportCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  showImportParams
	flagSet *pflag.FlagSet
}

type showImportParams struct {
	request.GetStorageImportDetailsRequest
}

func ShowImportCommand(service service.Storage) commands.Command {
	return &showImportCommand{
		BaseCommand: commands.New("show-import", "Show import task details"),
		service:     service,
	}
}

func (s *showImportCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if len(args) != 1 {
			return nil, fmt.Errorf("storage uuid is required")
		}

		s.params.UUID = args[0]
		var showImportStorageRequests []request.GetStorageImportDetailsRequest
		showImportStorageRequests = append(showImportStorageRequests, s.params.GetStorageImportDetailsRequest)

		var (
			mu                   sync.Mutex
			numOk                int
			storageImportDetails *upcloud.StorageImportDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			req := showImportStorageRequests[idx]
			msg := fmt.Sprintf("Show import task %q", req.UUID)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.GetStorageImportDetails(&req)
			storageImportDetails = details
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				e.SetMessage(fmt.Sprintf("%s: done", msg))
				e.SetDetails(details.UUID, "UUID: ")
				mu.Lock()
				numOk++
				mu.Unlock()
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(showImportStorageRequests),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(showImportStorageRequests) {
			return nil, fmt.Errorf("number of showImport operations that failed: %d", len(showImportStorageRequests)-numOk)
		}

		return storageImportDetails, nil
	}

}

func (s *showImportCommand) HandleOutput(writer io.Writer, out interface{}) error {

	details := out.(*upcloud.StorageImportDetails)

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "UUID: ", details.UUID)
	fmt.Fprintln(writer, "State: ", details.State)
	fmt.Fprintln(writer, "ClientContentLength: ", details.ClientContentLength)
	fmt.Fprintln(writer, "ClientContentType: ", details.ClientContentType)
	fmt.Fprintln(writer, "Completed: ", details.Completed)
	fmt.Fprintln(writer, "Created: ", details.Created)
	fmt.Fprintln(writer, "DirectUploadURL: ", details.DirectUploadURL)
	fmt.Fprintln(writer, "MD5Sum: ", details.MD5Sum)
	fmt.Fprintln(writer, "ReadBytes: ", details.ReadBytes)
	fmt.Fprintln(writer, "SHA256Sum: ", details.SHA256Sum)
	fmt.Fprintln(writer, "Source: ", details.Source)
	fmt.Fprintln(writer, "SourceLocation: ", details.SourceLocation)
	fmt.Fprintln(writer, "Error code: ", details.ErrorCode)
	fmt.Fprintln(writer, "Error message: ", details.ErrorMessage)
	fmt.Fprintln(writer)

	return nil
}
