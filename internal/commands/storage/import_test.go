package storage

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"os"
	"testing"
)

func TestImportCommand(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-")
	if err != nil {
		fmt.Println("Cannot create temporary file", err)
	}

	defer os.Remove(tmpFile.Name())

	var Storage1 = upcloud.Storage{
		UUID:   Uuid1,
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
		UUID:   Uuid2,
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
				"--source", upcloud.StorageImportSourceHTTPImport,
				"--source-location", "http://example.com",
			},
			request: request.CreateStorageImportRequest{
				StorageUUID:    Storage1.UUID,
				Source:         upcloud.StorageImportSourceHTTPImport,
				SourceLocation: "http://example.com",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedStorages = nil
			mss := mocks.MockStorageService{}
			mss.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)
			mss.On("CreateStorageImport", &test.request).Return(&StorageImportCompleted, nil)
			mss.On("GetStorageImportDetails", &request.GetStorageImportDetailsRequest{UUID: Storage1.UUID}).Return(&StorageImportCompleted, nil)
			mss.On("CreateStorage", mock.Anything).Return(&StorageDetails1, nil)

			ic := commands.BuildCommand(ImportCommand(&mss), nil, config.New(viper.New()))
			mocks.SetFlags(ic, test.args)

			_, err := ic.MakeExecuteCommand()(test.args)

			fmt.Println(err)
			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mss.AssertNumberOfCalls(t, "CreateStorageImport", 1)
				mss.AssertNumberOfCalls(t, "GetStorageImportDetails", 1)
				mss.AssertNumberOfCalls(t, "CreateStorage", 1)
			}
		})
	}
}
