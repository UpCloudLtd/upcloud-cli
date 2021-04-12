package storage

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListStorages(t *testing.T) {
	var Storage1 = upcloud.Storage{
		UUID:   UUID1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	var Storage2 = upcloud.Storage{
		UUID:   UUID2,
		Title:  Title2,
		Access: "private",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	var Storage3 = upcloud.Storage{
		UUID:   UUID3,
		Title:  Title3,
		Access: "public",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   10,
		Tier:   "maxiops",
	}
	var Storage4 = Storage1
	Storage4.Title = "mock-storage-title4"
	Storage4.Type = upcloud.StorageTypeCDROM
	var Storage5 = Storage1
	Storage5.Title = "mock-storage-title5"
	Storage5.Type = upcloud.StorageTypeTemplate
	var Storage6 = Storage3
	Storage6.Title = "mock-storage-title6"
	Storage6.Type = upcloud.StorageTypeBackup
	storages := upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2, Storage3, Storage4, Storage5, Storage6}}

	for _, testcase := range []struct {
		name    string
		private bool
		public  bool
		args    []string
		testFn  func(res upcloud.Storages, e error)
	}{
		{
			name: "List storages",
			args: []string{"--all"},
			testFn: func(res upcloud.Storages, e error) {
				assert.ElementsMatch(t, res.Storages, storages.Storages)
				assert.Nil(t, e)
			},
		},
		{
			name: "List public storages",
			args: []string{"--public"},
			testFn: func(res upcloud.Storages, e error) {
				assert.ElementsMatch(t, res.Storages, []upcloud.Storage{Storage3, Storage6})
				assert.Nil(t, e)
			},
		},
		{
			name: "List private by default",
			args: []string{},
			testFn: func(res upcloud.Storages, e error) {
				assert.ElementsMatch(t, res.Storages, []upcloud.Storage{Storage1, Storage2, Storage4, Storage5})
				assert.Nil(t, e)
			},
		},
		{
			name: "List cdrom",
			args: []string{"--cdrom"},
			testFn: func(res upcloud.Storages, e error) {
				assert.ElementsMatch(t, res.Storages, []upcloud.Storage{Storage4})
				assert.Nil(t, e)
			},
		},
		{
			name: "List public backup",
			args: []string{"--public", "--backup"},
			testFn: func(res upcloud.Storages, e error) {
				assert.ElementsMatch(t, res.Storages, []upcloud.Storage{Storage6})
				assert.Nil(t, e)
			},
		},
		{
			name: "List public template",
			args: []string{"--public", "--template"},
			testFn: func(res upcloud.Storages, e error) {
				assert.ElementsMatch(t, res.Storages, []upcloud.Storage{})
				assert.Nil(t, e)
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			CachedStorages = nil
			conf := config.New()
			mService := new(smock.Service)

			storages := upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2, Storage3, Storage4, Storage5, Storage6}}
			mService.On("GetStorages", mock.Anything).Return(&storages, nil)

			c := commands.BuildCommand(ListCommand(), nil, config.New())
			err := c.Cobra().Flags().Parse(testcase.args)
			assert.NoError(t, err)

			_, err = c.(commands.Command).Execute(commands.NewExecutor(conf, mService), "")
			assert.NoError(t, err)

			mService.AssertNumberOfCalls(t, "GetStorages", 1)
			// more checks
			// res, err := lc.MakeExecuteCommand()([]string{})
			// result := res.(*upcloud.Storages)
			// testcase.testFn(*result, err)
		})
	}
}
