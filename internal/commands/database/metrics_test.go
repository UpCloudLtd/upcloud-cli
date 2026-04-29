package database

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//go:embed metrics_testdata.json
var testdata []byte

func testMetricsCommand(t *testing.T, expected string, latestOnly bool) {
	t.Helper()
	text.DisableColors()

	metrics := upcloud.ManagedDatabaseMetrics{}
	err := json.Unmarshal(testdata, &metrics)
	if err != nil {
		t.Fatal(err)
	}

	mService := smock.Service{}
	mService.On("GetManagedDatabaseMetrics", mock.Anything).Return(&metrics, nil)

	conf := config.New()
	command := commands.BuildCommand(MetricsCommand(), nil, conf)

	args := []string{"test-uuid", "--period", "hour"}
	if latestOnly {
		args = append(args, "--latest-only")
	}

	command.Cobra().SetArgs(args)
	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}

func TestMetricsCommand(t *testing.T) {
	expected := `
  CPU usage %

     Time                  testdata-1 (master) 
    ───────────────────── ─────────────────────
     2026-04-29 13:47:30   6.58                
     2026-04-29 13:48:00   5.64                
    
  Memory usage %

     Time                  testdata-1 (master) 
    ───────────────────── ─────────────────────
     2026-04-29 13:47:30   65.43               
     2026-04-29 13:48:00   65.42               
    
  Disk space usage %

     Time                  testdata-1 (master) 
    ───────────────────── ─────────────────────
     2026-04-29 13:47:30   1.23                
     2026-04-29 13:48:00   1.23                
    
`
	testMetricsCommand(t, expected, false)
}

func TestMetricsCommand_LatestOnly(t *testing.T) {
	expected := `
  CPU usage %

     Time                  testdata-1 (master) 
    ───────────────────── ─────────────────────
     2026-04-29 13:48:00   5.64                
    
  Memory usage %

     Time                  testdata-1 (master) 
    ───────────────────── ─────────────────────
     2026-04-29 13:48:00   65.42               
    
  Disk space usage %

     Time                  testdata-1 (master) 
    ───────────────────── ─────────────────────
     2026-04-29 13:48:00   1.23                
    
`
	testMetricsCommand(t, expected, true)
}
