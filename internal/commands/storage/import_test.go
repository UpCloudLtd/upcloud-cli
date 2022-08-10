package storage

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/mockexecute"
	internal "github.com/UpCloudLtd/upcloud-cli/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
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

	fileNotFoundErr := "cannot get file size: stat testfile: no such file or directory"
	if runtime.GOOS == "windows" {
		fileNotFoundErr = "cannot get file size: CreateFile testfile: The system cannot find the file specified."
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
			name: "local import, non-existent file",
			args: []string{
				//				"--source-type", upcloud.StorageImportSourceDirectUpload,
				"--source-location", "testfile",
				"--zone", "fi-hel1",
				"--title", "test-2",
			},
			error: fileNotFoundErr,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedStorages = nil
			conf := config.New()
			mService := new(smock.Service)

			conf.Service = internal.Wrapper{Service: mService}

			mService.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)
			mService.On("CreateStorageImport", &test.request).Return(&StorageImportCompleted, nil)
			mService.On("GetStorageImportDetails", &request.GetStorageImportDetailsRequest{UUID: Storage1.UUID}).Return(&StorageImportCompleted, nil)
			mService.On("CreateStorage", mock.Anything).Return(&StorageDetails1, nil)

			c := commands.BuildCommand(ImportCommand(), nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, "CreateStorageImport", 1)
				mService.AssertNumberOfCalls(t, "GetStorageImportDetails", 1)
				mService.AssertNumberOfCalls(t, "CreateStorage", 1)
			}
		})
	}
}

func TestParseSource(t *testing.T) {
	for _, test := range []struct {
		name                    string
		input                   string
		expectedError           string
		expectedWindowsError    string
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
			name:                 "local non-existent file",
			input:                "foobar",
			expectedError:        "cannot get file size: stat foobar: no such file or directory",
			expectedWindowsError: "cannot get file size: CreateFile foobar: The system cannot find the file specified.",
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
				tempf, tempErr := ioutil.TempFile(t.TempDir(), "")
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
			if test.expectedError != "" {
				if test.expectedWindowsError != "" && runtime.GOOS == "windows" {
					assert.EqualError(t, err, test.expectedWindowsError)
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
