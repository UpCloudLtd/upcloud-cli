package storage

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/service"

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
	existingStorageUUIDOrName string
	noWait                    config.OptionalBoolean

	createParams createParams

	Resolver resolver.CachingStorage
}

// InitCommand implements Command.InitCommand
func (s *importCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.sourceLocation, "source-location", "", "Location of the source of the import. Can be a file or a URL.")
	flagSet.StringVar(&s.existingStorageUUIDOrName, "storage", "", "Import to an existing storage. Storage must be large enough and must be undetached or the server where the storage is attached must be in shutdown state.")
	config.AddToggleFlag(flagSet, &s.noWait, "no-wait", false, "Do not wait until the import finishes. Only applicable when importing from a remote URL.")
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
	// initial argument validation
	if s.sourceLocation == "" {
		return nil, fmt.Errorf("source-location required")
	}

	if s.existingStorageUUIDOrName == "" {
		if s.createParams.Zone == "" || s.createParams.Title == "" {
			return nil, fmt.Errorf("either existing storage or zone and title for a new storage to be created required")
		}
	} else if s.createParams.Zone != "" || s.createParams.Title != "" {
		return nil, fmt.Errorf("title and zone are not valid when using existing storage")
	}

	// figure out sourcetype and validate the inputs to the best of our ability
	parsedSource, sourceType, fileSize, err := parseSource(s.sourceLocation)
	if err != nil {
		return nil, err
	}

	// calculate filesize in gigabytes to validate storage sizes
	// add one because we're rounding down with integer division, otherwise we could end up consistently
	// creating too small storages to hold the file we want to upload
	fileSizeInGB := int(fileSize/1024/1024/1024) + 1

	// next, figure out if we want to import to an existing storage (and validate it) or create one
	var (
		storageToImportTo upcloud.Storage
	)
	if s.existingStorageUUIDOrName != "" {
		// user specified an existing storage, validate it
		// initialize resolver
		// TODO: maybe this resolver business should be rethought? this use case isnt really supported,
		//       possibly split resolving and caching to separate bits?
		_, err := s.Resolver.Get(exec.All())
		if err != nil {
			return nil, fmt.Errorf("cannot setup storage resolver: %w", err)
		}
		foundUUID, err := s.Resolver.Resolve(s.existingStorageUUIDOrName)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve existing storage: %w", err)
		}
		cached, err := s.Resolver.GetCached(foundUUID)
		if err != nil {
			return nil, fmt.Errorf("cannot get existing storage: %w", err)
		}
		if cached.Size < fileSizeInGB {
			return nil, fmt.Errorf("the existing storage is too small for the file")
		}
		storageToImportTo = cached
	} else {
		// We need to create a new storage.
		// Infer created storage size from the file if default size is used
		if s.createParams.Size == defaultCreateParams.Size && fileSizeInGB > defaultCreateParams.Size {
			s.createParams.Size = fileSizeInGB
		} else if s.createParams.Size < fileSizeInGB {
			// user gave a custom size, validate that it's large enough
			return nil, fmt.Errorf("the requested storage size is too small for the file")
		}
		createdStorage, err := createStorage(exec, &s.createParams)
		if err != nil {
			return nil, fmt.Errorf("cannot create storage: %w", err)
		}
		storageToImportTo = createdStorage
	}

	// input has been validated and we have a storage to import to, ready to start the actual import
	msg := fmt.Sprintf("importing to %v", storageToImportTo.UUID)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	startTime := time.Now()
	var (
		statusChan   = make(chan storageImportStatus)
		transferType string
	)
	switch sourceType {
	case upcloud.StorageImportSourceHTTPImport:
		// Import from the internet
		transferType = "download"
		result, err := exec.Storage().CreateStorageImport(
			&request.CreateStorageImportRequest{
				StorageUUID:    storageToImportTo.UUID,
				Source:         sourceType,
				SourceLocation: s.sourceLocation,
			})
		if err != nil {
			return nil, fmt.Errorf("failed to import: %w", err)
		}
		logline.SetMessage(fmt.Sprintf("%s: http import queued", msg))
		if !s.noWait.Value() {
			// start polling for import status if --no-wait was not entered on the command line
			go pollStorageImportStatus(exec.Storage(), storageToImportTo.UUID, statusChan)
		} else {
			// otherwise, we can just return the result and be done with it
			logline.SetMessage(fmt.Sprintf("%s: http import request sent", msg))
			logline.MarkDone()
			return output.OnlyMarshaled{Value: result}, nil
		}
	case upcloud.StorageImportSourceDirectUpload:
		// import from local file
		transferType = "upload"
		sourceFile, err := os.Open(parsedSource.Path)
		if err != nil {
			return nil, fmt.Errorf("cannot open local file: %w", err)
		}
		go importLocalFile(exec.Storage(), storageToImportTo.UUID, sourceFile, statusChan)
	}

	// import has been triggered, read updates from the process
	for statusUpdate := range statusChan {
		switch {
		case statusUpdate.err != nil:
			// we received an error, clean up log and return the error
			return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), statusUpdate.err)
		case statusUpdate.complete:
			// we're complete, clean up log and return the result
			logline.SetMessage(fmt.Sprintf("%s: done", msg))
			logline.MarkDone()
			return output.OnlyMarshaled{Value: statusUpdate.result}, nil
		case statusUpdate.bytesTransferred > 0:
			// got a status update
			bps := float64(statusUpdate.bytesTransferred) / time.Since(startTime).Seconds()
			// update the file size, if the backend returns a status update with it, eg. if
			// the import is a http import
			if fileSize == 0 && statusUpdate.result != nil {
				fileSize = int64(statusUpdate.result.ClientContentLength)
			}
			if fileSize > 0 {
				// we have knowledge of import file size, report progress percentage
				logline.SetMessage(fmt.Sprintf("%s: %sed %.2f%% (%sBps)",
					msg, transferType,
					float64(statusUpdate.bytesTransferred)/float64(fileSize)*100,
					ui.AbbrevNumBinaryPrefix(uint(bps)),
				))
			} else {
				// we have no knowledge of the remote file size, report bytes uploaded
				logline.SetMessage(fmt.Sprintf("%s: %sed %sB (%sBps)",
					msg, transferType,
					ui.AbbrevNumBinaryPrefix(uint(statusUpdate.bytesTransferred)),
					ui.AbbrevNumBinaryPrefix(uint(bps)),
				))
			}
		}
	}
	// status channel was closed but we did not receive either result or an error, fail.*/
	return nil, fmt.Errorf("upload aborted unexpectedly")
}

// TODO: figure out how to handle 'local http uploads', eg. piping from a local / non public internet url
//       if required(?)
func parseSource(location string) (parsedLocation *url.URL, sourceType string, fileSize int64, err error) {
	fileSize, err = getLocalFileSize(location)
	if err == nil {
		// we managed to open a local file with this path, so use that
		sourceType = upcloud.StorageImportSourceDirectUpload
		parsedLocation = &url.URL{Path: location}
		return
	}
	parsedLocation, err = url.Parse(location)
	switch {
	case err != nil:
		// error parsing url and can't open the file - return with error
		return nil, "", 0, fmt.Errorf("cannot parse url from source-location '%v': %w", location, err)
	case parsedLocation.Scheme == "" || parsedLocation.Scheme == "file":
		// parsed, but looks like a local file URL
		sourceType = upcloud.StorageImportSourceDirectUpload
		fileSize, err = getLocalFileSize(parsedLocation.Path)
		if err != nil {
			return nil, "", 0, fmt.Errorf("cannot get file size: %w", err)
		}
	default:
		// url was parsed and seems to not be a reference to a local file, make sure it's http or https
		sourceType = upcloud.StorageImportSourceHTTPImport
		if parsedLocation.Scheme != "http" && parsedLocation.Scheme != "https" {
			return nil, "", 0, fmt.Errorf("unsupported URL scheme '%v'", parsedLocation.Scheme)
		}
	}
	return
}

func createStorage(exec commands.Executor, params *createParams) (upcloud.Storage, error) {
	if err := params.processParams(); err != nil {
		return upcloud.Storage{}, err
	}
	msg := fmt.Sprintf("Creating storage %q", params.Title)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	details, err := exec.Storage().CreateStorage(&params.CreateStorageRequest)
	if err != nil {
		commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err) // nolint:errcheck
		return upcloud.Storage{}, err
	}
	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.SetDetails(details.UUID, "UUID: ")
	logline.MarkDone()
	return details.Storage, nil
}

func getLocalFileSize(path string) (size int64, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
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
