package storage

import (
	"bytes"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var storage1 = upcloud.Storage{
	UUID:   Uuid1,
	Title:  Title1,
	Access: "private",
	State:  "maintenance",
	Type:   "backup",
	Zone:   "fi-hel1",
	Size:   40,
	Tier:   "maxiops",
}

var storage2 = upcloud.Storage{
	UUID:   Uuid2,
	Title:  Title2,
	Access: "private",
	State:  "online",
	Type:   "normal",
	Zone:   "fi-hel1",
	Size:   40,
	Tier:   "maxiops",
}

var storage3 = upcloud.Storage{
	UUID:   Uuid3,
	Title:  Title3,
	Access: "public",
	State:  "online",
	Type:   "normal",
	Zone:   "fi-hel1",
	Size:   10,
	Tier:   "maxiops",
}

func TestListStorages(t *testing.T) {

	for _, testcase := range []struct {
		name    string
		private bool
		public  bool
		testFn  func(res upcloud.Storages, e error)
	}{
		{
			name:    "List storages",
			private: true,
			public:  true,
			testFn: func(res upcloud.Storages, e error) {
				assert.Equal(t, 2, len(res.Storages))
				assert.Nil(t, e)
			},
		},
		{
			name:    "List private storages",
			private: true,
			public:  false,
			testFn: func(res upcloud.Storages, e error) {
				assert.Equal(t, 2, len(res.Storages))
				assert.Nil(t, e)
			},
		},
		{
			name:    "List public storages",
			private: false,
			public:  true,
			testFn: func(res upcloud.Storages, e error) {
				assert.Equal(t, 2, len(res.Storages))
				assert.Nil(t, e)
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {

			mss := MockStorageService()
			mss.On("GetStorages", mock.Anything).Return(storages, nil)
			lc := commands.BuildCommand(ListCommand(mss), nil, config.New(viper.New()))

			res, err := lc.MakeExecuteCommand()([]string{})
			result := res.(*upcloud.Storages)
			testcase.testFn(*result, err)
		})
	}
}

func TestListStoragesOutput(t *testing.T) {
	storages := &upcloud.Storages{
		Storages: []upcloud.Storage{
			storage1,
			storage2,
			storage3,
		},
	}

	mss := MockStorageService()
	lc := commands.BuildCommand(ListCommand(mss), nil, config.New(viper.New()))

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
