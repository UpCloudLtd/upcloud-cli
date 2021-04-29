package storage

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/pflag"
)

// ImportCommand creates the "storage import" command
func ImportCommand() commands.Command {
	return &importCommand{
		BaseCommand: commands.New(
			"import",
			"Import a storage from external or local source",
			"upctl storage import --source-location https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/debian-10.9.0-amd64-netinst.iso --title my_storage --zone es-mad1",
		),
		createParams: newCreateParams(),
	}
}

type readerCounter struct {
	source io.Reader
	read   int64
}

// Read implements io.Reader
func (s *readerCounter) Read(p []byte) (n int, err error) {
	n, err = s.source.Read(p)
	atomic.AddInt64(&s.read, int64(n))
	return
}

func (s *readerCounter) counter() int {
	return int(atomic.LoadInt64(&s.read))
}

type importCommand struct {
	*commands.BaseCommand

	sourceLocation            string
	sourceType                string
	existingStorageUUIDOrName string
	wait                      config.OptionalBoolean

	createParams createParams

	Resolver resolver.CachingStorage
}

// InitCommand implements Command.InitCommand
func (s *importCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.sourceLocation, "source-location", "", "Location of the source of the import. Can be a file or a URL.")
	// TODO: is this flag actually required? Could we not just figure it out depending on the location?
	// eg. if there's a file with the name given in location, use direct upload, otherwise validate url and use http import?
	flagSet.StringVar(&s.sourceType, "source-type", "", fmt.Sprintf("Source type, is derived from source-location if not given. Available: %s,%s",
		upcloud.StorageImportSourceHTTPImport,
		upcloud.StorageImportSourceDirectUpload))
	flagSet.StringVar(&s.existingStorageUUIDOrName, "storage", "", "Import to an existing storage. Storage must be large enough and must be undetached or the server where the storage is attached must be in shutdown state.")
	config.AddToggleFlag(flagSet, &s.wait, "wait", true, fmt.Sprintf("Wait until the import finishes. Implied if source is set to %s",
		upcloud.StorageImportSourceDirectUpload))
	applyCreateFlags(flagSet, &s.createParams, defaultCreateParams)
	s.AddFlags(flagSet)
}

type storageImportStatus struct {
	result           *upcloud.StorageImportDetails
	bytesTransferred int64
	err              error
	// we need separate cmoplete as the local and remote imports report in a different manner
	// with remote import polling the details and returning a new result every time
	complete bool
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *importCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {

	svc := exec.Storage()

	var (
		existingStorage *upcloud.Storage
		storageUUID     string
	)

	if s.existingStorageUUIDOrName != "" {
		// initialize resolver
		// TODO: maybe this should be rethought?
		_, err := s.Resolver.Get(exec.All())
		if err != nil {
			return nil, fmt.Errorf("cannot setup storage resolver: %w", err)
		}
		if s.sourceLocation == "" {
			return nil, errors.New("source-location must be defined")
		}
		foundUUID, err := s.Resolver.Resolve(s.existingStorageUUIDOrName)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve existing storage: %w", err)
		}
		cached, err := s.Resolver.GetCached(foundUUID)
		if err != nil {
			return nil, fmt.Errorf("cannot get existing storage: %w", err)
		}
		existingStorage = &cached
		storageUUID = cached.UUID
	} else if s.sourceLocation == "" || s.createParams.Zone == "" || s.createParams.Title == "" {
		return nil, errors.New("source-location and either existing storage or both zone and title are required")
	}

	// Infer source type from source location
	// TODO: is there any sense in passing this as a parameter?
	// TODO: is there any sense that this is not the domain of the sdk?
	if s.sourceType == "" {
		parsedURL, err := url.Parse(s.sourceLocation)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Scheme == "file" {
			s.sourceType = upcloud.StorageImportSourceDirectUpload
		} else {
			s.sourceType = upcloud.StorageImportSourceHTTPImport
		}
	}
	var (
		sourceFile    *os.File
		localFileSize int64
	)
	if s.sourceType == upcloud.StorageImportSourceDirectUpload {
		// this could be done in the actual upload call, but we'd rather validate all input in the first place
		f, err := os.Open(s.sourceLocation)
		if err != nil {
			return nil, err
		}
		sourceFile = f
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}
		localFileSize = stat.Size()
		fileSizeInGB := int(localFileSize / 1024 / 1024 / 1024)
		if existingStorage != nil && existingStorage.Size < fileSizeInGB {
			return nil, fmt.Errorf("the existing storage is too small for the file")
		}
		if existingStorage == nil && s.createParams.Size != defaultCreateParams.Size &&
			s.createParams.Size < fileSizeInGB {
			return nil, fmt.Errorf("the requested storage size is too small for the file")
		}
		// Infer created storage size from the file if default size is used
		if existingStorage == nil && s.createParams.Size == defaultCreateParams.Size &&
			fileSizeInGB > defaultCreateParams.Size {
			s.createParams.Size = fileSizeInGB
		}
	}
	if s.sourceType == upcloud.StorageImportSourceHTTPImport {
		_, err := url.Parse(s.sourceLocation)
		if err != nil {
			return nil, fmt.Errorf("invalid import url: %w", err)
		}
	}

	// Create storage, if it doesn't exist yet
	if existingStorage == nil {
		if err := s.createParams.processParams(); err != nil {
			return nil, err
		}
		msg := fmt.Sprintf("Creating storage %q", s.createParams.Title)
		logline := exec.NewLogEntry(msg)
		logline.StartedNow()
		details, err := svc.CreateStorage(&s.createParams.CreateStorageRequest)
		if err != nil {
			logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
			logline.SetDetails(err.Error(), "error: ")
			return nil, fmt.Errorf("import failed: %w", err)
		}
		logline.SetMessage(fmt.Sprintf("%s: done", msg))
		logline.SetDetails(details.UUID, "UUID: ")
		logline.MarkDone()
		existingStorage = &details.Storage
		storageUUID = details.UUID
	}

	// Create import task
	msg := fmt.Sprintf("importing to %v", existingStorage.UUID)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	startTime := time.Now()
	var (
		statusChan   = make(chan storageImportStatus)
		transferType string
	)
	if sourceFile != nil {
		// Import from local file
		transferType = "upload"
		go importLocalFile(svc, storageUUID, sourceFile, statusChan)
	} else {
		// Import from http location
		transferType = "download"
		result, err := svc.CreateStorageImport(
			&request.CreateStorageImportRequest{
				StorageUUID:    storageUUID,
				Source:         upcloud.StorageImportSourceHTTPImport,
				SourceLocation: s.sourceLocation,
			})
		if err != nil {
			return nil, fmt.Errorf("failed to import: %w", err)
		}
		logline.SetMessage(fmt.Sprintf("%s: http import queued", msg))
		if s.wait {
			// start polling for import status if --wait was entered
			go pollStorageImportStatus(svc, storageUUID, statusChan)
		} else {
			// otherwise, we can just return the result
			return output.OnlyMarshaled{Value: result}, nil
		}
	}

	// wait for updates from the import process
	for statusUpdate := range statusChan {
		switch {
		case statusUpdate.err != nil:
			// we received an error, clean up log and return the error
			logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, statusUpdate.err.Error()))
			logline.SetDetails(statusUpdate.err.Error(), "error: ")
			return nil, statusUpdate.err
		case statusUpdate.complete:
			// we're complete, clean up log and return the result
			logline.SetMessage(fmt.Sprintf("%s: done", msg))
			logline.MarkDone()
			return output.OnlyMarshaled{Value: statusUpdate.result}, nil
		case statusUpdate.bytesTransferred > 0:
			// got a status update
			bps := float64(statusUpdate.bytesTransferred) / time.Since(startTime).Seconds()
			// get the file size, if possible - clientContentLength can still be 0
			importFileSize := localFileSize
			if importFileSize == 0 && statusUpdate.result != nil {
				importFileSize = int64(statusUpdate.result.ClientContentLength)
			}
			if importFileSize > 0 {
				// we have knowledge of import file size, report progress percentage
				logline.SetMessage(fmt.Sprintf("%s: %sed %.2f%% (%sbps)",
					msg, transferType,
					float64(statusUpdate.bytesTransferred)/float64(importFileSize)*100,
					ui.AbbrevNum(uint(bps)),
				))
			} else {
				// we have no knowledge of the remote file size, report bytes uploaded
				logline.SetMessage(fmt.Sprintf("%s: %sed %sB (%sBps)",
					msg, transferType,
					ui.AbbrevNum(uint(statusUpdate.bytesTransferred)),
					ui.AbbrevNum(uint(bps)),
				))
			}
		}
	}
	// status channel was closed but we did not receive either result or an error, fail.
	return nil, fmt.Errorf("upload aborted unexpectedly")
}

func pollStorageImportStatus(svc service.Storage, uuid string, statusChan chan<- storageImportStatus) {
	// make sure we close the channel when exiting poller
	defer close(statusChan)

	sleepSecs := 2
	for {
		details, err := svc.GetStorageImportDetails(&request.GetStorageImportDetailsRequest{
			UUID: uuid,
		})
		switch {
		case err != nil:
			statusChan <- storageImportStatus{err: err}
			return
		case details.ErrorCode != "":
			statusChan <- storageImportStatus{err: fmt.Errorf("%s (%s)", details.ErrorMessage, details.ErrorCode)}
			return
		case details.State == upcloud.StorageImportStateCancelled:
			statusChan <- storageImportStatus{err: fmt.Errorf("cancelled")}
			return
		case details.State == upcloud.StorageImportStateCompleted:
			statusChan <- storageImportStatus{result: details, complete: true}
			return
		}
		if read := details.ReadBytes; read > 0 {
			statusChan <- storageImportStatus{result: details, bytesTransferred: int64(read)}
		}
		time.Sleep(time.Duration(sleepSecs) * time.Second)
	}
}

func importLocalFile(svc service.Storage, uuid string, file *os.File, statusChan chan<- storageImportStatus) {
	// make sure we close the channel when exiting import
	defer close(statusChan)
	chDone := make(chan storageImportStatus)
	reader := &readerCounter{source: file}

	// figure out content type
	contentType := "application/octet-stream"
	switch filepath.Ext(file.Name()) {
	case ".gz":
		contentType = "application/gzip"
	case ".xz":
		contentType = "application/x-xz"
	}

	go func() {
		imported, err := svc.CreateStorageImport(
			&request.CreateStorageImportRequest{
				StorageUUID:    uuid,
				ContentType:    contentType,
				Source:         upcloud.StorageImportSourceDirectUpload,
				SourceLocation: reader,
			})
		chDone <- storageImportStatus{result: imported, err: err, complete: true}
	}()
	updateTicker := time.NewTicker(300 * time.Millisecond)
	for {
		select {
		case result := <-chDone:
			statusChan <- result
			return
		case <-updateTicker.C:
			statusChan <- storageImportStatus{bytesTransferred: int64(reader.counter())}
		}
	}
}
