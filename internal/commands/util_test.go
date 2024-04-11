package commands_test

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"

	"github.com/stretchr/testify/assert"
)

func TestParseSSHKeys(t *testing.T) {
	type args struct {
		sshKeys []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr assert.ErrorAssertionFunc
	}{
		{name: "string", args: struct{ sshKeys []string }{sshKeys: []string{`ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMWq/xsiYPgA/HLsaWHcjAGnwU+pJy9BUmvIlMBpkdn2 noname`}}, want: []string{`ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMWq/xsiYPgA/HLsaWHcjAGnwU+pJy9BUmvIlMBpkdn2 noname`}, wantErr: func(_ assert.TestingT, err error, _ ...interface{}) bool {
			return err == nil
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := commands.ParseSSHKeys(tt.args.sshKeys)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseSSHKeys(%v)", tt.args.sshKeys)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseSSHKeys(%v)", tt.args.sshKeys)
		})
	}
}
