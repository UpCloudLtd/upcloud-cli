package storage

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStorageHumanOutput(t *testing.T) {
	s := upcloud.Storage{
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
		Storage: s,
		BackupRule: &upcloud.BackupRule{
			Interval:  "daily",
			Time:      "0400",
			Retention: 7,
		},
	}

	expected :=
		`  Common │        UUID │ 01101f27-196f-47e9-a055-4e2e8bb3b419  
         │       Title │ test-1                                
         │        Zone │ fi-hel1                               
         │       State │ online                                
         │  Size (GiB) │ 10                                    
         │        Type │ normal                                
         │        Tier │ maxiops                               
         │     License │ 0                                     
         │     Created │                                       
         │      Origin │                                       
─────────┼─────────────────────────────────────────────────────
 Servers │ no servers using this storage                       
─────────┼─────────────────────────────────────────────────────
  Backup │  Backup Rule │   Interval │ daily                   
         │              │       Time │ 0400                    
         │              │  Retention │ 7                       `

	output, _ := ShowCommand(nil).HandleOutput(sd)
	assert.Equal(t, expected, output)
}
