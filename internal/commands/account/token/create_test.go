package token

import (
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateToken(t *testing.T) {
	created := time.Now()

	for _, test := range []struct {
		name  string
		resp  *upcloud.Token
		args  []string
		req   request.CreateTokenRequest
		errFn assert.ErrorAssertionFunc
	}{
		{
			name: "defaults",
			args: []string{
				"--name", "test",
				"--expires-in", "1h",
			},
			req: request.CreateTokenRequest{
				Name:               "test",
				ExpiresAt:          created.Add(1 * time.Hour),
				CanCreateSubTokens: false,
				AllowedIPRanges:    nil,
			},
			resp: &upcloud.Token{
				APIToken:           "ucat_01JH5D3ZZJVZS6JC713FA11CB8",
				ID:                 "0cd8eab4-ecb7-445b-a457-6019b0a00496",
				Name:               "test",
				Type:               "workspace",
				CreatedAt:          created,
				ExpiresAt:          created.Add(1 * time.Hour),
				LastUsed:           nil,
				CanCreateSubTokens: false,
				AllowedIPRanges:    []string{"0.0.0.0/0", "::/0"},
			},
			errFn: assert.NoError,
		},
		{
			name: "missing name",
			args: []string{
				"--expires-in", "1h",
			},
			errFn: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, `required flag(s) "name" not set`)
			},
		},
		{
			name: "invalid expires-in",
			args: []string{
				"--name", "test",
				"--expires-in", "seppo",
			},
			errFn: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, `invalid argument "seppo" for "--expires-in"`)
			},
		},
		{
			name: "invalid expires-at",
			args: []string{
				"--name", "test",
				"--expires-at", "seppo",
			},
			errFn: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, `invalid expires-at: `)
			},
		},
		{
			name: "missing expiry",
			args: []string{
				"--name", "test",
			},
			errFn: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, `either expires-in or expires-at must be set`)
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := CreateCommand()
			mService := new(smock.Service)

			if test.resp != nil {
				mService.On("CreateToken", mock.MatchedBy(func(req *request.CreateTokenRequest) bool {
					// service uses time.Now() with "expires-in" added to it to set ExpiresAt, so we can't set a mock to any
					// static value. Instead, we'll just check that the request has the correct name and that the ExpiresAt
					// is within 1 second of "now".
					return assert.Equal(t, test.req.Name, req.Name) && assert.InDelta(t, test.req.ExpiresAt.UnixMilli(), req.ExpiresAt.UnixMilli(), 1000)
				})).Once().Return(test.resp, nil)
			}

			c := commands.BuildCommand(testCmd, nil, conf)
			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.errFn(t, err) {
				mService.AssertExpectations(t)
			}
		})
	}
}
