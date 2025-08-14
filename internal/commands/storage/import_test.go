package storage

import (
	"io"
	"io/fs"
	"net/url"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReaderCounterInterface(_ *testing.T) {
	rc := &readerCounter{}
	var _ io.Reader = rc
}

func TestImportCommand(t *testing.T) {
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
		name               string
		args               []string
		error              string
		expectFileNotExist bool
		request            request.CreateStorageImportRequest
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
			error: `required flag(s) "source-location" not set`,
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
			name: "use existing storage",
			args: []string{
				"--source-location", "http://example.com",
				"--storage", Storage2.Title,
			},
			request: request.CreateStorageImportRequest{
				StorageUUID:    Storage2.UUID,
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
			expectFileNotExist: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedStorages = nil
			conf := config.New()
			mService := new(smock.Service)

			mService.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)
			req := test.request
			mService.On("CreateStorageImport", &req).Return(&StorageImportCompleted, nil)
			mService.On("GetStorageImportDetails", &request.GetStorageImportDetailsRequest{UUID: req.StorageUUID}).Return(&StorageImportCompleted, nil)
			mService.On("CreateStorage", mock.Anything).Return(&StorageDetails1, nil)

			c := commands.BuildCommand(ImportCommand(), nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.error != "" || test.expectFileNotExist {
				assert.Error(t, err)
				if test.expectFileNotExist {
					assert.ErrorIs(t, err, fs.ErrNotExist)
				} else {
					assert.EqualError(t, err, test.error)
				}
			} else {
				assert.NoError(t, err)

				createCount := 1
				if slices.Contains(test.args, "--storage") {
					createCount = 0
				}

				mService.AssertNumberOfCalls(t, "CreateStorageImport", 1)
				mService.AssertNumberOfCalls(t, "GetStorageImportDetails", 1)
				mService.AssertNumberOfCalls(t, "CreateStorage", createCount)
			}
		})
	}
}

func TestParseSource(t *testing.T) {
	for _, test := range []struct {
		name                    string
		input                   string
		expectedError           string
		expectFileNotExist      bool
		expectedSourceType      string
		expectedFileSize        int64
		expectedParsedURLScheme string
		expectedParsedURLHost   string
		expectedParsedURLPath   string
	}{
		{
			name:                  "raw local filename",
			input:                 "tempfile",
			expectedSourceType:    upcloud.StorageImportSourceDirectUpload,
			expectedFileSize:      5,
			expectedParsedURLPath: "tempfile",
		},
		{
			name:                  "local file with file:// path",
			input:                 "file://tempfile",
			expectedSourceType:    upcloud.StorageImportSourceDirectUpload,
			expectedFileSize:      5,
			expectedParsedURLPath: "tempfile",
		},
		{
			name:               "local non-existent file",
			input:              "foobar",
			expectFileNotExist: true,
		},
		{
			name:                    "remote http url",
			input:                   "http://127.0.0.1/remotefile",
			expectedSourceType:      upcloud.StorageImportSourceHTTPImport,
			expectedFileSize:        0,
			expectedParsedURLScheme: "http",
			expectedParsedURLHost:   "127.0.0.1",
			expectedParsedURLPath:   "/remotefile",
		},
		{
			name:                    "remote https url",
			input:                   "https://127.0.0.1/remotefile",
			expectedSourceType:      upcloud.StorageImportSourceHTTPImport,
			expectedFileSize:        0,
			expectedParsedURLScheme: "https",
			expectedParsedURLHost:   "127.0.0.1",
			expectedParsedURLPath:   "/remotefile",
		},
		{
			name:          "remote ftp url",
			input:         "ftp://127.0.0.1/remotefile",
			expectedError: "unsupported URL scheme 'ftp'",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var (
				purl  *url.URL
				stype string
				fsize int64
				err   error
			)
			testFileName := test.input
			if strings.Contains(testFileName, "tempfile") {
				tempf, tempErr := os.CreateTemp(t.TempDir(), "")
				assert.NoError(t, tempErr)
				_, tempErr = tempf.WriteString("hello")
				assert.NoError(t, tempErr)
				tempErr = tempf.Close()
				assert.NoError(t, tempErr)
				testFileName = tempf.Name()
				purl, stype, fsize, err = parseSource(strings.Replace(testFileName, "tempfile", tempf.Name(), 1))
			} else {
				purl, stype, fsize, err = parseSource(test.input)
			}
			if test.expectedError != "" || test.expectFileNotExist {
				assert.Error(t, err)
				if test.expectFileNotExist {
					assert.ErrorIs(t, err, fs.ErrNotExist)
				} else {
					assert.EqualError(t, err, test.expectedError)
				}
				assert.Nil(t, purl)
				assert.Equal(t, stype, "")
				assert.Equal(t, fsize, int64(0))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, purl)
				assert.Equal(t, purl.Scheme, test.expectedParsedURLScheme)
				assert.Equal(t, purl.Host, test.expectedParsedURLHost)
				if test.expectedParsedURLPath == "tempfile" {
					assert.Equal(t, purl.Path, testFileName)
				} else {
					assert.Equal(t, purl.Path, test.expectedParsedURLPath)
				}
				assert.Equal(t, stype, test.expectedSourceType)
				assert.Equal(t, fsize, test.expectedFileSize)
			}
		})
	}
}
