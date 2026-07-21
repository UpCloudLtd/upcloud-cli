package database

import (
	"strings"
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

func TestDatabasePlans_SortedHumanOutput(t *testing.T) {
	text.DisableColors()
	plans := []upcloud.ManagedDatabasePlanServiceType{
		{
			Type: "pg",
			ComputeShapes: []upcloud.ManagedDatabasePlanComputeShape{
				{
					Compute:  "regular-shape",
					Family:   "regular",
					CPU:      2,
					MemoryGB: 8,
					NodeCounts: []int{
						1, 2,
					},
					Backups: []string{"regular", "extended"},
					Storage: upcloud.ManagedDatabasePlanStorageInfo{
						Options: []upcloud.ManagedDatabasePlanStorageOption{
							{BaseGiB: 50, MaxGiB: 50},
							{BaseGiB: 100, MaxGiB: 200},
						},
					},
				},
			},
		},
	}

	legacyPg := upcloud.ManagedDatabaseType{
		ServicePlans: []upcloud.ManagedDatabaseServicePlan{
			{Plan: "test-plan-5", NodeCount: 3, CoreNumber: 16, MemoryAmount: 128, StorageSize: 2048, BackupConfig: upcloud.ManagedDatabaseBackupConfig{MaxCount: 30}},
			{Plan: "test-plan-3", NodeCount: 1, CoreNumber: 2, MemoryAmount: 16, StorageSize: 256, BackupConfig: upcloud.ManagedDatabaseBackupConfig{MaxCount: 15}},
			{Plan: "test-plan-1", NodeCount: 1, CoreNumber: 1, MemoryAmount: 8, StorageSize: 2048, BackupConfig: upcloud.ManagedDatabaseBackupConfig{MaxCount: 3}},
			{Plan: "test-plan-2", NodeCount: 1, CoreNumber: 2, MemoryAmount: 4, StorageSize: 2048, BackupConfig: upcloud.ManagedDatabaseBackupConfig{MaxCount: 7}},
			{Plan: "test-plan-4", NodeCount: 3, CoreNumber: 16, MemoryAmount: 128, StorageSize: 1024, BackupConfig: upcloud.ManagedDatabaseBackupConfig{MaxCount: 20}},
		},
	}

	mService := smock.Service{}
	mService.On("GetManagedDatabasePlans", mock.Anything).Return(plans, nil)
	mService.On("GetManagedDatabaseServiceType", mock.Anything).Return(&legacyPg, nil)

	conf := config.New()
	// force human output
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(PlansCommand(), nil, conf)

	command.Cobra().SetArgs([]string{"pg", "--show-legacy"})
	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.NoError(t, err)
	assert.Contains(t, output, "regular, extended")
	assert.NotContains(t, output, ", regular, extended")
	assert.Contains(t, output, "1, 2")
	assert.NotContains(t, output, ", 1, 2")
	assert.Contains(t, output, "50, 100-200")
	assert.NotContains(t, output, ", 50, 100-200")
	assert.Contains(t, output, "3 PITR days")
	assert.Contains(t, output, "legacy")
	assert.Less(t, strings.Index(output, "test-plan-1"), strings.Index(output, "test-plan-2"))
	assert.Less(t, strings.Index(output, "test-plan-2"), strings.Index(output, "test-plan-3"))
	assert.Less(t, strings.Index(output, "test-plan-3"), strings.Index(output, "test-plan-4"))
	assert.Less(t, strings.Index(output, "test-plan-4"), strings.Index(output, "test-plan-5"))
}
