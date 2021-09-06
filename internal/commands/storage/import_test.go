package storage_test

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/storage"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	internal "github.com/UpCloudLtd/upcloud-cli/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestImportCommand(t *testing.T) {
	t.Parallel()
	tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-")
	if err != nil {
		t.Fatalf("Cannot create temporary file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})

	Storage1 := upcloud.Storage{
		UUID:   UUID1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	StorageDetails1 := upcloud.StorageDetails{
		Storage: Storage1,
		BackupRule: &upcloud.BackupRule{
			Interval:  "sun",
			Time:      "0800",
			Retention: 5,
		},
	}
	Storage2 := upcloud.Storage{
		UUID:   UUID2,
		Title:  Title2,
		Access: "private",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}

	StorageImportCompleted := upcloud.StorageImportDetails{
		State: upcloud.StorageImportStateCompleted,
	}

	for _, test := range []struct {
		name         string
		args         []string
		error        string
		request      request.CreateStorageImportRequest
		windowsError string
	}{
		{
			name: "source is missing",
			args: []string{
				"--source-location", "http://example.com",
				"--zone", "fi-hel1",
				"--title", "test-1",
			},
			request: request.CreateStorageImportRequest{
				StorageUUID:    Storage1.UUID,
				Source:         upcloud.StorageImportSourceHTTPImport,
				SourceLocation: "http://example.com",
			},
		},
		{
			name: "location is missing",
			args: []string{
				//				"--source-type", upcloud.StorageImportSourceHTTPImport,
				"--zone", "fi-hel1",
				"--title", "test-1",
			},
			error: "source-location required",
		},
		{
			name: "http import",
			args: []string{
				//				"--source-type", upcloud.StorageImportSourceHTTPImport,
				"--source-location", "http://example.com",
				"--zone", "fi-hel1",
				"--title", "test-2",
			},
			request: request.CreateStorageImportRequest{
				StorageUUID:    Storage1.UUID,
				Source:         upcloud.StorageImportSourceHTTPImport,
				SourceLocation: "http://example.com",
			},
		},
		{
			name: "local import, non-existent file",
			args: []string{
				//				"--source-type", upcloud.StorageImportSourceDirectUpload,
				"--source-location", "testfile",
				"--zone", "fi-hel1",
				"--title", "test-2",
			},
			error:        "cannot get file size: stat testfile: no such file or directory",
			windowsError: "cannot get file size: CreateFile testfile: The system cannot find the file specified.",
		},
	} {
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			storage.CachedStorages = nil
			conf := config.New()
			mService := new(smock.Service)

			conf.Service = internal.Wrapper{Service: mService}

			mService.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)
			mService.On("CreateStorageImport", &test.request).Return(&StorageImportCompleted, nil)
			mService.On("GetStorageImportDetails", &request.GetStorageImportDetailsRequest{UUID: Storage1.UUID}).Return(&StorageImportCompleted, nil)
			mService.On("CreateStorage", mock.Anything).Return(&StorageDetails1, nil)

			c := commands.BuildCommand(storage.ImportCommand(), nil, conf)
			err := c.Cobra().Flags().Parse(test.args)
			assert.NoError(t, err)

			_, err = c.(commands.NoArgumentCommand).ExecuteWithoutArguments(commands.NewExecutor(conf, mService, flume.New("test")))

			if test.error != "" {
				if test.windowsError != "" && runtime.GOOS == "windows" {
					assert.EqualError(t, err, test.windowsError)
				} else {
					assert.EqualError(t, err, test.error)
				}
			} else {
				mService.AssertNumberOfCalls(t, "CreateStorageImport", 1)
				mService.AssertNumberOfCalls(t, "GetStorageImportDetails", 1)
				mService.AssertNumberOfCalls(t, "CreateStorage", 1)
			}
		})
	}
}
