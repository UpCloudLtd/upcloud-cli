package storage

import (
	"io"
	"io/ioutil"
	"net/url"
	"runtime"
	"strings"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestReaderCounterInterface(t *testing.T) {
	t.Parallel()
	rc := &readerCounter{}
	var _ io.Reader = rc
}

func TestParseSource(t *testing.T) {
	t.Parallel()
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
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
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
