package storage

import (
  "github.com/UpCloudLtd/cli/internal/commands"
  "github.com/UpCloudLtd/cli/internal/config"
  "github.com/UpCloudLtd/cli/internal/mocks"
  "github.com/spf13/viper"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "testing"
)

func TestRestoreBackupCommand(t *testing.T) {
  methodName := "RestoreBackup"


  for _, test := range []struct {
    name string
    args []string
    methodCalls int
  }{
    {
      name: "Backend called",
      args: []string{},
      methodCalls: 1,
    },
  }{
    t.Run(test.name, func(t *testing.T) {
      mss := MockStorageService()
      mss.On(methodName, mock.Anything).Return(nil, nil)

      tc := commands.BuildCommand(RestoreBackupCommand(mss), nil, config.New(viper.New()))
      mocks.SetFlags(tc, test.args)

      results, err := tc.MakeExecuteCommand()([]string{Storage2.UUID})
      for _, result := range results.([]interface{}) {
        assert.Equal(t, nil, result)
      }

      assert.Nil(t, err)

      mss.AssertNumberOfCalls(t, methodName, test.methodCalls)
    })
  }
}
