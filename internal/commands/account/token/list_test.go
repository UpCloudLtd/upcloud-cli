package token

import (
	"bytes"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/testutils"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/gemalto/flume"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

func TestListCommand(t *testing.T) {
	text.DisableColors()

	tokens := upcloud.Tokens{
		{
			ID:                 "0cdabbf9-090b-4fc5-a6ae-3f76801ed171",
			Name:               "test",
			Type:               "workspace",
			CreatedAt:          *testutils.MustParseRFC3339(t, "2021-10-01T12:00:00Z"),
			ExpiresAt:          *testutils.MustParseRFC3339(t, "2022-10-01T12:00:00Z"),
			LastUsed:           testutils.MustParseRFC3339(t, "2021-11-01T12:00:00Z"),
			CanCreateSubTokens: false,
			AllowedIPRanges:    []string{"0.0.0.0/0", "::/0"},
		},
		{
			ID:                 "0cdabbf9-090b-4fc5-a6ae-3f76801ed172",
			Name:               "foo",
			Type:               "workspace",
			CreatedAt:          *testutils.MustParseRFC3339(t, "2021-10-01T12:00:00Z"),
			ExpiresAt:          *testutils.MustParseRFC3339(t, "2022-10-01T12:00:00Z"),
			LastUsed:           nil,
			CanCreateSubTokens: true,
			AllowedIPRanges:    []string{"0.0.0.0/0", "::/0"},
		},
	}

	expected := `
 UUID                                   Name   Type        Last Used                       Expires At                    
────────────────────────────────────── ────── ─────────── ─────────────────────────────── ───────────────────────────────
 0cdabbf9-090b-4fc5-a6ae-3f76801ed171   test   workspace   2021-11-01 12:00:00 +0000 UTC   2022-10-01 12:00:00 +0000 UTC 
 0cdabbf9-090b-4fc5-a6ae-3f76801ed172   foo    workspace   <nil>                           2022-10-01 12:00:00 +0000 UTC 

`

	svc := &smock.Service{}
	conf := config.New()

	svc.On("GetTokens").Once().Return(&tokens, nil)
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(ListCommand(), nil, conf)
	out, err := command.(commands.NoArgumentCommand).ExecuteWithoutArguments(commands.NewExecutor(conf, svc, flume.New("test")))
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, conf.Output(), out)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())

	svc.AssertExpectations(t)
}
