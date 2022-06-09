package database

import (
	"bytes"
	"strings"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/gemalto/flume"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDatabasePlans_SortedHumanOutput(t *testing.T) {
	text.DisableColors()
	plansPg := upcloud.ManagedDatabaseType{
		ServicePlans: []upcloud.ManagedDatabaseServicePlan{
			{Plan: "test-plan-5", NodeCount: 3, CoreNumber: 16, MemoryAmount: 128, StorageSize: 2048},
			{Plan: "test-plan-3", NodeCount: 1, CoreNumber: 2, MemoryAmount: 16, StorageSize: 256},
			{Plan: "test-plan-1", NodeCount: 1, CoreNumber: 1, MemoryAmount: 8, StorageSize: 2048},
			{Plan: "test-plan-2", NodeCount: 1, CoreNumber: 2, MemoryAmount: 4, StorageSize: 2048},
			{Plan: "test-plan-4", NodeCount: 3, CoreNumber: 16, MemoryAmount: 128, StorageSize: 1024},
		},
	}

	mService := smock.Service{}
	mService.On("GetManagedDatabaseServiceType", mock.Anything).Return(&plansPg, nil)

	conf := config.New()
	// force human output
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(PlansCommand(), nil, conf)

	res, err := command.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, &mService, flume.New("test")), "pg")

	assert.Nil(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, conf, res)
	output := buf.String()

	assert.NoError(t, err)

	assert.Less(t, strings.Index(output, "test-plan-1"), strings.Index(output, "test-plan-2"))
	assert.Less(t, strings.Index(output, "test-plan-2"), strings.Index(output, "test-plan-3"))
	assert.Less(t, strings.Index(output, "test-plan-3"), strings.Index(output, "test-plan-4"))
	assert.Less(t, strings.Index(output, "test-plan-4"), strings.Index(output, "test-plan-5"))
}
