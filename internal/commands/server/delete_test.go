package server

import (
  "github.com/UpCloudLtd/cli/internal/commands"
  "github.com/UpCloudLtd/cli/internal/config"
  "github.com/UpCloudLtd/cli/internal/mocks"
  "github.com/spf13/viper"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "testing"
)

func TestDeleteServerCommand(t *testing.T) {
  deleteServer := "DeleteServer"
  deleteServerAndStorages := "DeleteServerAndStorages"

  for _, test := range []struct {
    name string
    args []string
    deleteServCalls int
    deleteServStorageCalls int
  }{
    {
      name: "Backend called, no args",
      args: []string{},
      deleteServCalls: 0,
      deleteServStorageCalls: 1,
    },
    {
      name: "Delete-storages true",
      args: []string{"--delete-storages", "true"},
      deleteServCalls: 0,
      deleteServStorageCalls: 1,
    },
    {
      name: "Delete-storages false",
      args: []string{"--delete-storages", "false"},
      deleteServCalls: 1,
      deleteServStorageCalls: 0,
    },
  }{
    t.Run(test.name, func(t *testing.T) {
      mss := MockServerService()
      mss.On(deleteServer, mock.Anything).Return(nil, nil)
      mss.On(deleteServerAndStorages, mock.Anything).Return(nil, nil)

      tc := commands.BuildCommand(DeleteCommand(mss), nil, config.New(viper.New()))
      mocks.SetFlags(tc, test.args)

      results, err := tc.MakeExecuteCommand()([]string{Server1.UUID})
      for _, result := range results.([]interface{}) {
        assert.Nil(t, result)
      }

      assert.Nil(t, err)

      mss.AssertNumberOfCalls(t, deleteServer, test.deleteServCalls)
      mss.AssertNumberOfCalls(t, deleteServerAndStorages, test.deleteServStorageCalls)
    })
  }
}
