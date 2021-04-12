package storage

import (
	"bytes"
	"testing"
	"time"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"
	"github.com/UpCloudLtd/cli/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestStorageHumanOutput(t *testing.T) {
	storage := upcloud.Storage{
		Access:     "private",
		License:    0,
		PartOfPlan: "",
		Size:       10,
		State:      "online",
		Tier:       "maxiops",
		Title:      "test-1",
		Type:       "normal",
		UUID:       "01101f27-196f-47e9-a055-4e2e8bb3b419",
		Zone:       "fi-hel1",
		Origin:     "",
		Created:    time.Time{},
	}

	sd := &upcloud.StorageDetails{
		Storage: storage,
		BackupRule: &upcloud.BackupRule{
			Interval:  "daily",
			Time:      "0400",
			Retention: 7,
		},
		ServerUUIDs: []string{"0077fa3d-32db-4b09-9f5f-30d9e9afb565"},
		BackupUUIDs: []string{"012580a1-32a1-466e-a323-689ca16f2d43"},
	}

	storages := []upcloud.Storage{
		storage,
	}

	// sid := &upcloud.StorageImportDetails{
	// 	ClientContentLength: 1,
	// 	ClientContentType:   "abc",
	// 	ErrorCode:           "ghi",
	// 	ErrorMessage:        "jkl",
	// 	MD5Sum:              "mno",
	// 	ReadBytes:           2,
	// 	SHA256Sum:           "pqr",
	// 	Source:              "directUpload",
	// 	State:               "prepared",
	// 	UUID:                "07a6c9a3-300e-4d0e-b935-624f3dbdff3f",
	// 	WrittenBytes:        3,
	// }

	expected := `  
  Storage
    UUID:    01101f27-196f-47e9-a055-4e2e8bb3b419 
    Title:   test-1                               
    type:    normal                               
    State:   online                               
    Size:    10                                   
    Tier:    maxiops                              
    Zone:    fi-hel1                              
    Server:  0077fa3d-32db-4b09-9f5f-30d9e9afb565 
    Origin:                                       
    Created: 0001-01-01 00:00:00 +0000 UTC        
    Licence: 0                                    

  
  Backup Rule
    Interval:  daily 
    Time:      0400  
    Retention: 7     

  Available Backups

     UUID                                 
    ──────────────────────────────────────
     012580a1-32a1-466e-a323-689ca16f2d43 
    
`

	mService := smock.Service{}
	mService.On("GetStorageDetails", &request.GetStorageDetailsRequest{UUID: storage.UUID}).Return(sd, nil)
	mService.On("GetStorages", &request.GetStoragesRequest{}).Return(&upcloud.Storages{Storages: storages}, nil)

	conf := config.New()
	// force human output
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(ShowCommand(), nil, conf)

	// get resolver to initialize command cache
	_, err := command.(*showCommand).Get(&mService)
	if err != nil {
		t.Fatal(err)
	}
	res, err := command.(commands.Command).Execute(commands.NewExecutor(conf, &mService), storage.UUID)

	assert.Nil(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, conf, res)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())
}
