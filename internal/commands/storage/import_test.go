package storage

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"
	internal "github.com/UpCloudLtd/cli/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReaderCounterInterface(t *testing.T) {
	rc := &readerCounter{}
	var _ io.Reader = rc
}

func TestImportCommand(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-")
	if err != nil {
		fmt.Println("Cannot create temporary file", err)
	}

	defer os.Remove(tmpFile.Name())

	var Storage1 = upcloud.Storage{
		UUID:   UUID1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	var StorageDetails1 = upcloud.StorageDetails{
		Storage: Storage1,
		BackupRule: &upcloud.BackupRule{
			Interval:  "sun",
			Time:      "0800",
			Retention: 5,
		},
	}
	var Storage2 = upcloud.Storage{
		UUID:   UUID2,
		Title:  Title2,
		Access: "private",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}

	var StorageImportCompleted = upcloud.StorageImportDetails{
		State: upcloud.StorageImportStateCompleted,
	}

	for _, test := range []struct {
		name    string
		args    []string
		error   string
		request request.CreateStorageImportRequest
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
			name: "http import",
			args: []string{
				"--source-type", upcloud.StorageImportSourceHTTPImport,
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
	} {
		t.Run(test.name, func(t *testing.T) {
			// TODO: remove this
			t.SkipNow()
			CachedStorages = nil
			conf := config.New()
			mService := new(smock.Service)

			conf.Service = internal.Wrapper{Service: mService}

			mService.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)
			mService.On("CreateStorageImport", &test.request).Return(&StorageImportCompleted, nil)
			mService.On("GetStorageImportDetails", &request.GetStorageImportDetailsRequest{UUID: Storage1.UUID}).Return(&StorageImportCompleted, nil)
			mService.On("CreateStorage", mock.Anything).Return(&StorageDetails1, nil)

			c := commands.BuildCommand(ImportCommand(), nil, conf)
			err := c.Cobra().Flags().Parse(test.args)
			assert.NoError(t, err)

			_, err = c.(commands.Command).Execute(commands.NewExecutor(conf, mService), "")

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, "CreateStorageImport", 1)
				mService.AssertNumberOfCalls(t, "GetStorageImportDetails", 1)
				mService.AssertNumberOfCalls(t, "CreateStorage", 1)
			}
		})
	}
}
