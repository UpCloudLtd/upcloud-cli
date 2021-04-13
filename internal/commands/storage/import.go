package storage

import (
	"errors"
	"fmt"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

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
			"",
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
	sourceFile                *os.File
	sourceFileSize            int
	existingStorageUUIDOrName string
	wait                      bool

	storageUUID  string
	createParams createParams

	Resolver resolver.CachingStorage
}

// InitCommand implements Command.InitCommand
func (s *importCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.sourceLocation, "source-location", "", "Location of the source of the import. Can be a file or a URL.")
	flagSet.StringVar(&s.sourceType, "source-type", "", fmt.Sprintf("Source type, is derived from source-location if not given. Available: %s,%s",
		upcloud.StorageImportSourceHTTPImport,
		upcloud.StorageImportSourceDirectUpload))
	flagSet.StringVar(&s.existingStorageUUIDOrName, "storage", "", "Import to an existing storage. Storage must be large enough and must be undetached or the server where the storage is attached must be in shutdown state.")
	flagSet.BoolVar(&s.wait, "wait", true, fmt.Sprintf("Wait until the import finishes. Implied if source is set to %s",
		upcloud.StorageImportSourceDirectUpload))
	applyCreateFlags(flagSet, &s.createParams, defaultCreateParams)
	s.AddFlags(flagSet)
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
	var contentType string
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
		s.sourceFile = f
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}
		s.sourceFileSize = int(stat.Size())
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

	// Create storage
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
	var createdStorageImport *upcloud.StorageImportDetails
	if existingStorage != nil {
		msg := fmt.Sprintf("importing to storage %v", existingStorage.UUID)
		logline := exec.NewLogEntry(msg)
		logline.StartedNow()
		var err error
		// Import from local file
		if s.sourceFile != nil {
			chDone := make(chan struct{})
			var importErr error
			reader := &readerCounter{source: s.sourceFile}
			go func() {
				createdStorageImport, importErr = svc.CreateStorageImport(
					&request.CreateStorageImportRequest{
						StorageUUID:    s.storageUUID,
						ContentType:    contentType,
						Source:         s.sourceType,
						SourceLocation: reader,
					})
				chDone <- struct{}{}
			}()
			var prevRead int
			sleepSecs := 2
			sleepTicker := time.NewTicker(time.Duration(sleepSecs) * time.Second)
		loop:
			for {
				select {
				case <-chDone:
					if importErr == nil {
						logline.SetMessage(fmt.Sprintf("%s: done", msg))
						break loop
					}
					return nil, fmt.Errorf("failed to import: %w", err)
				case <-sleepTicker.C:
					if read := reader.counter(); read > 0 {
						logline.SetMessage(fmt.Sprintf("%s: uploaded %.2f%% (%sbps)",
							msg,
							float64(read)/float64(s.sourceFileSize)*100,
							ui.AbbrevNum(uint(read-prevRead)*8/uint(sleepSecs)),
						))
						prevRead = read
					}
				}
			}
		} else {
			// Import from http source
			createdStorageImport, err = svc.CreateStorageImport(
				&request.CreateStorageImportRequest{
					StorageUUID:    s.storageUUID,
					ContentType:    contentType,
					Source:         s.sourceType,
					SourceLocation: s.sourceLocation,
				})
			if err != nil {
				return nil, fmt.Errorf("failed to import: %w", err)
			}
			logline.SetMessage(fmt.Sprintf("%s: http import queued", msg))
			if s.wait {
				var prevRead int
				sleepSecs := 5
				for {
					details, importErr := svc.GetStorageImportDetails(&request.GetStorageImportDetailsRequest{
						UUID: existingStorage.UUID,
					})
					switch {
					case importErr != nil:
						return nil, fmt.Errorf("can not get details: %w", importErr)
					case details.ErrorCode != "":
						return nil, fmt.Errorf("%s (%s)", details.ErrorMessage, details.ErrorCode)
					case details.State == upcloud.StorageImportStateCancelled:
						return nil, fmt.Errorf("%s: cancelled", msg)
					case details.State == upcloud.StorageImportStateCompleted:
						logline.SetMessage(fmt.Sprintf("%s: done", msg))
						logline.MarkDone()
						return output.OnlyMarshaled{Value: createdStorageImport}, nil
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
			} else {
				logline.SetMessage(fmt.Sprintf("%s: http import request sent", msg))
				logline.MarkDone()
			}
		}
	}
	return nil, fmt.Errorf("no storage to import to")
}
