package database

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseListTitleFallback(t *testing.T) {
	text.DisableColors()
	databases := []upcloud.ManagedDatabase{
		{UUID: "091f1afe-4ddd-4d43-afad-6aa3069cc7fe", Title: "service-name", Name: "hostname-prefix-1", State: "running"},
		{UUID: "091f1afe-4ddd-4d43-afad-6aa3069cc7fe", Name: "hostname-prefix-2", State: "running"},
	}

	for _, test := range []struct {
		name string
		args []string
		page request.Page
	}{
		{
			name: "default page",
			page: request.Page{Size: 100, Number: 0},
		},
		{
			name: "limit and page args",
			page: request.Page{Size: 18, Number: 19},
			args: []string{"--limit", "18", "--page", "19"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			page := test.page

			mService := smock.Service{}
			mService.On("GetManagedDatabases", &request.GetManagedDatabasesRequest{Page: &page}).Return(databases, nil)

			conf := config.New()
			command := commands.BuildCommand(ListCommand(), nil, conf)
			command.Cobra().SetArgs(test.args)

			output, err := mockexecute.MockExecute(command, &mService, conf)

			assert.NoError(t, err)
			assert.Regexp(t, `UUID\s+Title\s+Type\s+Plan\s+Zone\s+State`, output)
			assert.Contains(t, output, "service-name")
			assert.NotContains(t, output, "hostname-prefix-1")
			assert.Contains(t, output, "hostname-prefix-2")
		})
	}
}
