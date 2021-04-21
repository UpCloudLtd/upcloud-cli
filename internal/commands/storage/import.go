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

	storageUUID  string
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

type storageImportResult struct {
	result *upcloud.StorageImportDetails
	err    error
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *importCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {

	svc := exec.Storage()

	var existingStorage *upcloud.Storage
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
		s.storageUUID = cached.UUID
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
		contentType string
		sourceFile  *os.File
	)
	if s.sourceType == upcloud.StorageImportSourceDirectUpload {
		f, err := os.Open(s.sourceLocation)
		if err != nil {
			return nil, err
		}
		switch filepath.Ext(s.sourceLocation) {
		case ".gz":
			contentType = "application/gzip"
		case ".xz":
			contentType = "application/x-xz"
		default:
			contentType = "application/octet-stream"
		}
		sourceFile = f
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}
		fileSizeInGB := int(stat.Size() / 1024 / 1024 / 1024)
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

	if existingStorage == nil {
		if err := s.createParams.processParams(); err != nil {
			return nil, err
		}
	}

	var (
		createdStorage *upcloud.StorageDetails
	)

	// Create storage, if it doesn't exist yet
	if existingStorage == nil {
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
		createdStorage = details
		existingStorage = &details.Storage
		s.storageUUID = createdStorage.UUID
	}

	// Create import task
	msg := fmt.Sprintf("importing to storage %v", existingStorage.UUID)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	// Import from local file
	var (
		result *upcloud.StorageImportDetails
		err    error
	)
	if sourceFile != nil {
		result, err = importLocalFile(svc, logline, msg, s.storageUUID, contentType, sourceFile)
	} else {
		result, err = importRemoteFile(svc, logline, msg, s.storageUUID, contentType, s.sourceLocation, s.wait)
	}
	if err != nil {
		return nil, err
	}
	return output.OnlyMarshaled{Value: result}, nil
}

func importRemoteFile(svc service.Storage, logline *ui.LogEntry, msg string, uuid string, contentType string, location string, wait bool) (*upcloud.StorageImportDetails, error) {
	createdStorageImport, err := svc.CreateStorageImport(
		&request.CreateStorageImportRequest{
			StorageUUID:    uuid,
			ContentType:    contentType,
			Source:         upcloud.StorageImportSourceHTTPImport,
			SourceLocation: location,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to import: %w", err)
	}
	logline.SetMessage(fmt.Sprintf("%s: http import queued", msg))
	if wait {
		var prevRead int
		sleepSecs := 5
		for {
			details, err := svc.GetStorageImportDetails(&request.GetStorageImportDetailsRequest{
				UUID: uuid,
			})
			switch {
			case err != nil:
				return nil, fmt.Errorf("can not get details: %w", err)
			case details.ErrorCode != "":
				return nil, fmt.Errorf("%s (%s)", details.ErrorMessage, details.ErrorCode)
			case details.State == upcloud.StorageImportStateCancelled:
				return nil, fmt.Errorf("%s: cancelled", msg)
			case details.State == upcloud.StorageImportStateCompleted:
				logline.SetMessage(fmt.Sprintf("%s: done", msg))
				logline.MarkDone()
				return createdStorageImport, nil
			}
			if read := details.ReadBytes; read > 0 {
				if details.ClientContentLength > 0 {
					logline.SetMessage(fmt.Sprintf("%s: downloaded %.2f%% (%sbps)",
						msg,
						float64(read)/float64(details.ClientContentLength)*100,
						ui.AbbrevNum(uint(read-prevRead)*8/uint(sleepSecs)),
					))
					prevRead = read
				} else {
					logline.SetMessage(fmt.Sprintf("%s: downloaded %s (%sbps)",
						msg,
						ui.FormatBytes(read),
						ui.AbbrevNum(uint(read-prevRead)*8/uint(sleepSecs)),
					))
				}
			}
			time.Sleep(time.Duration(sleepSecs) * time.Second)
		}
	}
	logline.SetMessage(fmt.Sprintf("%s: http import request sent", msg))
	logline.MarkDone()
	return createdStorageImport, nil
}

func importLocalFile(svc service.Storage, logline *ui.LogEntry, logmsg, uuid, contentType string, file *os.File) (*upcloud.StorageImportDetails, error) {
	chDone := make(chan storageImportResult)
	reader := &readerCounter{source: file}
	fileSize := int64(0)
	if stat, err := file.Stat(); err == nil {
		fileSize = stat.Size()
	} else {
		return nil, fmt.Errorf("cannot stat input file: %w", err)
	}
	go func() {
		imported, err := svc.CreateStorageImport(
			&request.CreateStorageImportRequest{
				StorageUUID:    uuid,
				ContentType:    contentType,
				Source:         upcloud.StorageImportSourceDirectUpload,
				SourceLocation: reader,
			})
		chDone <- storageImportResult{imported, err}
	}()
	var prevRead int
	sleepSecs := 2
	sleepTicker := time.NewTicker(time.Duration(sleepSecs) * time.Second)
	for {
		select {
		case result := <-chDone:
			if result.err == nil {
				logline.SetMessage(fmt.Sprintf("%s: done", logmsg))
				return result.result, nil
			}
			return nil, fmt.Errorf("failed to import: %w", result.err)
		case <-sleepTicker.C:
			if read := reader.counter(); read > 0 {
				logline.SetMessage(fmt.Sprintf("%s: uploaded %.2f%% (%sbps)",
					logmsg,
					float64(read)/float64(fileSize)*100,
					ui.AbbrevNum(uint(read-prevRead)*8/uint(sleepSecs)),
				))
				prevRead = read
			}
		}
	}
}
