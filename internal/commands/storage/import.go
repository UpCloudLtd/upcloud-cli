package storage

import (
	"errors"
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"io"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
)

// ImportCommand creates the "storage import" command
func ImportCommand(service service.Storage) commands.Command {
	return &importCommand{
		BaseCommand: commands.New("import", "Import a storage from external or local source"),
		service:     service,
	}
}

var defaultImportParams = &importParams{
	CreateStorageImportRequest: request.CreateStorageImportRequest{},
	createStorage:              newCreateParams(),
	wait:                       true,
}

func newImportParams() importParams {
	return importParams{
		CreateStorageImportRequest: request.CreateStorageImportRequest{},
		createStorage:              newCreateParams(),
	}
}

type importParams struct {
	request.CreateStorageImportRequest
	createStorage             createParams
	sourceLocation            string
	existingStorageUUIDOrName string
	wait                      bool

	sourceFile      *os.File
	sourceFileSize  int
	existingStorage *upcloud.Storage
}

func (s *importParams) processParams(srv service.Storage) error {
	if s.existingStorageUUIDOrName != "" {
		storage, err := searchStorage(&CachedStorages, srv, s.existingStorageUUIDOrName, true)
		if err != nil {
			return err
		}
		s.existingStorage = storage[0]
		s.CreateStorageImportRequest.StorageUUID = storage[0].UUID
	}
	if s.sourceLocation == "" || s.createStorage.Zone == "" || s.createStorage.Title == "" {
		return errors.New("source-location, zone and title are required")
	}
	// Infer source type from source location
	if s.Source == "" {
		parsedURL, err := url.Parse(s.sourceLocation)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Scheme == "file" {
			s.Source = upcloud.StorageImportSourceDirectUpload
		} else {
			s.Source = upcloud.StorageImportSourceHTTPImport
		}
	}
	if s.Source == upcloud.StorageImportSourceDirectUpload {
		f, err := os.Open(s.sourceLocation)
		if err != nil {
			return err
		}
		switch filepath.Ext(s.sourceLocation) {
		case ".gz":
			s.ContentType = "application/gzip"
		case ".xz":
			s.ContentType = "application/x-xz"
		default:
			s.ContentType = "application/octet-stream"
		}
		s.sourceFile = f
		stat, err := f.Stat()
		if err != nil {
			return err
		}
		s.sourceFileSize = int(stat.Size())
		if s.existingStorage != nil && float64(s.existingStorage.Size) < float64(stat.Size()/1024/1024/1024) {
			return fmt.Errorf("the existing storage is too small for the file")
		}
		if s.existingStorage == nil && s.createStorage.Size != defaultCreateParams.Size &&
			float64(s.createStorage.Size) < float64(stat.Size()/1024/1024/1024) {
			return fmt.Errorf("the requested storage size is too small for the file")
		}
		// Infer created storage size from the file if default size is used
		if s.existingStorage == nil && s.createStorage.Size == defaultCreateParams.Size &&
			float64(stat.Size()/1024/1024/1024) > float64(defaultCreateParams.Size) {
			s.createStorage.Size = int(math.Ceil(float64(stat.Size() / 1024 / 1024 / 1024)))
		}
	}
	if s.Source == upcloud.StorageImportSourceHTTPImport {
		_, err := url.Parse(s.sourceLocation)
		if err != nil {
			return err
		}
		s.SourceLocation = s.sourceLocation
	}
	return nil
}

func (s *importParams) close() {
	if s.sourceFile != nil {
		_ = s.sourceFile.Close()
	}
}

type readerCounter struct {
	source io.Reader
	read   int64
}

func (s *readerCounter) readCounting(p []byte) (n int, err error) {
	n, err = s.source.Read(p)
	atomic.AddInt64(&s.read, int64(n))
	return
}

func (s *readerCounter) counter() int {
	return int(atomic.LoadInt64(&s.read))
}

type importCommand struct {
	*commands.BaseCommand
	service      service.Storage
	importParams importParams
	flagSet      *pflag.FlagSet
}

func importFlags(fs *pflag.FlagSet, dst, def *importParams) {
	fs.StringVar(&dst.sourceLocation, "source-location", def.sourceLocation, "Location of the source of the import. Can be a file or a URL.\n[Required]")
	fs.StringVar(&dst.Source, "source", def.Source, fmt.Sprintf("Source type. Available: %s,%s",
		upcloud.StorageImportSourceHTTPImport,
		upcloud.StorageImportSourceDirectUpload))
	fs.StringVar(&dst.existingStorageUUIDOrName, "storage", def.existingStorageUUIDOrName, "Import to an existing storage. Storage must be large enough and must be undetached "+"or the server where the storage is attached must be in shutdown state.")
	fs.BoolVar(&dst.wait, "wait", def.wait, fmt.Sprintf("Wait until the import finishes. Implied if source is set to %s",
		upcloud.StorageImportSourceDirectUpload))
	createFlags(fs, &dst.createStorage, defaultCreateParams)
}

// InitCommand implements Command.InitCommand
func (s *importCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getStorageArgumentCompletionFunction(s.service))
	s.flagSet = &pflag.FlagSet{}
	s.importParams = newImportParams()
	importFlags(s.flagSet, &s.importParams, defaultImportParams)
	s.AddFlags(s.flagSet)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *importCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		errorOrGenericError := func(err error) error {
			if s.Config().InteractiveUI() {
				return errors.New("import failed")
			}
			return err
		}
		if err := s.importParams.processParams(s.service); err != nil {
			return nil, err
		}
		if s.importParams.existingStorage == nil {
			if err := s.importParams.createStorage.processParams(); err != nil {
				return nil, err
			}
		}

		var (
			createdStorage *upcloud.StorageDetails
			workFlowErr    error
		)

		// Create storage
		handlerCreateStorage := func(idx int, e *ui.LogEntry) {
			msg := fmt.Sprintf("Creating storage %q", s.importParams.createStorage.Title)
			e.SetMessage(msg)
			e.StartedNow()
			details, err := s.service.CreateStorage(&s.importParams.createStorage.CreateStorageRequest)
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
				workFlowErr = err
			} else {
				e.SetMessage(fmt.Sprintf("%s: done", msg))
				e.SetDetails(details.UUID, "UUID: ")
				createdStorage = details
				s.importParams.existingStorage = &details.Storage
			}
		}
		if s.importParams.existingStorage == nil {
			ui.StartWorkQueue(ui.WorkQueueConfig{
				NumTasks:           1,
				MaxConcurrentTasks: 1,
				EnableUI:           s.Config().InteractiveUI(),
			}, handlerCreateStorage)
			if workFlowErr != nil {
				return nil, errorOrGenericError(workFlowErr)
			}
			s.importParams.CreateStorageImportRequest.StorageUUID = createdStorage.UUID
		}

		// Create import task
		var createdStorageImport *upcloud.StorageImportDetails
		handlerImport := func(idx int, e *ui.LogEntry) {
			msg := "Importing to storage"
			e.SetMessage(msg)
			e.StartedNow()
			var err error
			// Import from local file
			if s.importParams.sourceFile != nil {
				chDone := make(chan struct{})
				var importErr error
				reader := &readerCounter{source: s.importParams.sourceFile}
				s.importParams.SourceLocation = reader
				go func() {
					createdStorageImport, importErr = s.service.CreateStorageImport(
						&s.importParams.CreateStorageImportRequest)
					chDone <- struct{}{}
				}()
				wait := true
				var prevRead int
				sleepSecs := 2
				for {
					select {
					case <-chDone:
						err = importErr
						wait = false
					default:
					}
					if read := reader.counter(); read > 0 {
						e.SetMessage(fmt.Sprintf("%s: uploaded %.2f%% (%sbps)",
							msg,
							float64(read)/float64(s.importParams.sourceFileSize)*100,
							ui.AbbrevNum(uint(read-prevRead)*8/uint(sleepSecs)),
						))
						prevRead = read
					}
					if !wait {
						if err == nil {
							e.SetMessage(fmt.Sprintf("%s: done", msg))
						}
						goto end
					}
					time.Sleep(time.Duration(sleepSecs) * time.Second)
				}
			}
			// Import from http source
			if s.importParams.sourceFile == nil {
				createdStorageImport, err = s.service.CreateStorageImport(
					&s.importParams.CreateStorageImportRequest)
				if err != nil {
					goto end
				}
				e.SetMessage(fmt.Sprintf("%s: http import queued", msg))
				if s.importParams.wait {
					var prevRead int
					sleepSecs := 5
					for {
						details, importErr := s.service.GetStorageImportDetails(&request.GetStorageImportDetailsRequest{
							UUID: s.importParams.existingStorage.UUID,
						})
						switch {
						case importErr != nil:
							err = importErr
							goto end
						case details.ErrorCode != "":
							err = fmt.Errorf("%s (%s)", details.ErrorMessage, details.ErrorCode)
							goto end
						case details.State == upcloud.StorageImportStateCancelled:
							err = fmt.Errorf("%s: cancelled", msg)
							goto end
						case details.State == upcloud.StorageImportStateCompleted:
							e.SetMessage(fmt.Sprintf("%s: done", msg))
							goto end
						}
						if importErr != nil {
						}
						if details.ErrorCode != "" {
						}
						if read := details.ReadBytes; read > 0 {
							if details.ClientContentLength > 0 {
								e.SetMessage(fmt.Sprintf("%s: downloaded %.2f%% (%sbps)",
									msg,
									float64(read)/float64(details.ClientContentLength)*100,
									ui.AbbrevNum(uint(read-prevRead)*8/uint(sleepSecs)),
								))
								prevRead = read
							} else {
								e.SetMessage(fmt.Sprintf("%s: downloaded %s (%sbps)",
									msg,
									ui.FormatBytes(read),
									ui.AbbrevNum(uint(read-prevRead)*8/uint(sleepSecs)),
								))
							}
						}
						time.Sleep(time.Duration(sleepSecs) * time.Second)
					}
				}
			}

		end:
			if err != nil {
				e.SetMessage(ui.DefaultErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
				workFlowErr = err
			}
		}
		if s.importParams.existingStorage != nil {
			ui.StartWorkQueue(ui.WorkQueueConfig{
				NumTasks:           1,
				MaxConcurrentTasks: 1,
				EnableUI:           s.Config().InteractiveUI(),
			}, handlerImport)
			if workFlowErr != nil {
				return nil, errorOrGenericError(workFlowErr)
			}
		}

		return map[string]interface{}{
			"created_storage": createdStorage,
			"import_task":     createdStorageImport,
		}, nil
	}
}
