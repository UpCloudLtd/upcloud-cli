package storage

import (
	"bytes"
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListStorages(t *testing.T) {

	var Storage1 = upcloud.Storage{
		UUID:   Uuid1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	var Storage2 = upcloud.Storage{
		UUID:   Uuid2,
		Title:  Title2,
		Access: "private",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	var Storage3 = upcloud.Storage{
		UUID:   Uuid3,
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
			args: []string{"--private", "--public"},
			testFn: func(res upcloud.Storages, e error) {
				assert.ElementsMatch(t, res.Storages, storages.Storages)
				assert.Nil(t, e)
			},
		},
		{
			name: "List private storages",
			args: []string{"--private"},
			testFn: func(res upcloud.Storages, e error) {
				assert.ElementsMatch(t, res.Storages, []upcloud.Storage{Storage1, Storage2, Storage4, Storage5})
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
			stgs := storages
			mss := new(mocks.MockStorageService)
			mss.On("GetStorages", mock.Anything).Return(&stgs, nil)

			lc := commands.BuildCommand(ListCommand(mss), nil, config.New(viper.New()))
			_ = mocks.SetFlags(lc, testcase.args)
			res, err := lc.MakeExecuteCommand()([]string{})
			result := res.(*upcloud.Storages)
			testcase.testFn(*result, err)
		})
	}
}

func TestListStoragesOutput(t *testing.T) {
	Title1 = "mock-storage-title1"
	Title2 = "mock-storage-title2"
	Title3 = "mock-storage-title3"
	Uuid1 = "0127dfd6-3884-4079-a948-3a8881df1a7a"
	Uuid2 = "012bde1d-f0e7-4bb2-9f4a-74e1f2b49c07"
	Uuid3 = "012c61a6-b8f0-48c2-a63a-b4bf7d26a655"

	var Storage1 = upcloud.Storage{
		UUID:   Uuid1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}

	var Storage2 = upcloud.Storage{
		UUID:   Uuid2,
		Title:  Title2,
		Access: "private",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}

	var Storage3 = upcloud.Storage{
		UUID:   Uuid3,
		Title:  Title3,
		Access: "public",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   10,
		Tier:   "maxiops",
	}
	var storages = &upcloud.Storages{
		Storages: []upcloud.Storage{
			Storage1,
			Storage2,
			Storage3,
		},
	}

	lc := commands.BuildCommand(ListCommand(new(mocks.MockStorageService)), nil, config.New(viper.New()))

	expected := `
 UUID                                   Title                 Zone      State         Type     Size   Tier      Created 
────────────────────────────────────── ───────────────────── ───────── ───────────── ──────── ────── ───────── ─────────
 0127dfd6-3884-4079-a948-3a8881df1a7a   mock-storage-title1   fi-hel1   maintenance   backup     40   maxiops           
 012bde1d-f0e7-4bb2-9f4a-74e1f2b49c07   mock-storage-title2   fi-hel1   online        normal     40   maxiops           
 012c61a6-b8f0-48c2-a63a-b4bf7d26a655   mock-storage-title3   fi-hel1   online        normal     10   maxiops           

`

	buf := new(bytes.Buffer)
	err := lc.HandleOutput(buf, storages)

	assert.Nil(t, err)
	assert.Equal(t, expected, buf.String())
}
