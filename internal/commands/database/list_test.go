package database

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v7/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDatabaseListTitleFallback(t *testing.T) {
	text.DisableColors()
	databases := []upcloud.ManagedDatabase{
		{UUID: "091f1afe-4ddd-4d43-afad-6aa3069cc7fe", Title: "service-name", Name: "hotname-prefix-1", State: "running"},
		{UUID: "091f1afe-4ddd-4d43-afad-6aa3069cc7fe", Name: "hotname-prefix-2", State: "running"},
	}

	mService := smock.Service{}
	mService.On("GetManagedDatabases", mock.Anything).Return(databases, nil)

	conf := config.New()
	command := commands.BuildCommand(ListCommand(), nil, conf)

	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.NoError(t, err)
	assert.Regexp(t, `UUID\s+Title\s+Type\s+Plan\s+Zone\s+State`, output)
	assert.Contains(t, output, "service-name")
	assert.NotContains(t, output, "hotname-prefix-1")
	assert.Contains(t, output, "hotname-prefix-2")
}
