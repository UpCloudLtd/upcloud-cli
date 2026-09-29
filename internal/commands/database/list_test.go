package database

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseListTitleFallback(t *testing.T) {
	text.DisableColors()
	databases := `[
		{"uuid": "091f1afe-4ddd-4d43-afad-6aa3069cc7fe", "title": "service-name", "name": "hostname-prefix-1", "state": "running"},
		{"uuid": "091f1afe-4ddd-4d43-afad-6aa3069cc7fe", "name": "hostname-prefix-2", "state": "running"}
	]`

	for _, test := range []struct {
		name   string
		args   []string
		limit  string
		offset string
	}{
		{
			name:   "default page",
			limit:  "100",
			offset: "0",
		},
		{
			name:   "limit and page args",
			args:   []string{"--limit", "18", "--page", "19"},
			limit:  "18",
			offset: "324",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			api := newFakeAPI(t, map[string]fakeResponse{
				"GET /1.3/database": {status: http.StatusOK, body: databases},
			})

			conf := config.New()
			command := commands.BuildCommand(ListCommand(), nil, conf)
			command.Cobra().SetArgs(test.args)

			output, err := mockexecute.MockExecuteWithV9(command, &smock.Service{}, api.client, conf)

			require.NoError(t, err)
			assert.Regexp(t, `UUID\s+Title\s+Type\s+Plan\s+Zone\s+State`, output)
			assert.Contains(t, output, "service-name")
			assert.NotContains(t, output, "hostname-prefix-1")
			assert.Contains(t, output, "hostname-prefix-2")

			requests := api.requestsTo(http.MethodGet, "/1.3/database")
			require.Len(t, requests, 1)
			assert.Equal(t, test.limit, requests[0].Query.Get("limit"))
			assert.Equal(t, test.offset, requests[0].Query.Get("offset"))
		})
	}
}

func TestDatabaseListMachineReadableOutputUsesV9Response(t *testing.T) {
	const id = "091f1afe-4ddd-4d43-afad-6aa3069cc7fe"
	databases := `[{"uuid":"` + id + `","title":"service-name","type":"mysql","plan_components":{"compute":{"name":"rdb.standard.2CPU-8GB","node_count":2}}}]`

	for _, outputFormat := range []string{config.ValueOutputJSON, config.ValueOutputYAML} {
		t.Run(outputFormat, func(t *testing.T) {
			api := newFakeAPI(t, map[string]fakeResponse{"GET /1.3/database": {status: http.StatusOK, body: databases}})
			conf := config.New()
			conf.Viper().Set(config.KeyOutput, outputFormat)
			command := commands.BuildCommand(ListCommand(), nil, conf)

			out, err := mockexecute.MockExecuteWithV9(command, &smock.Service{}, api.client, conf)
			require.NoError(t, err)
			assert.Contains(t, out, "plan_components")
			assert.Contains(t, out, "rdb.standard.2CPU-8GB")
			if outputFormat == config.ValueOutputJSON {
				var got []map[string]any
				require.NoError(t, json.Unmarshal([]byte(out), &got))
				require.Len(t, got, 1)
				assert.Equal(t, id, got[0]["uuid"])
				assert.Contains(t, got[0], "plan_components")
			}
		})
	}
}
