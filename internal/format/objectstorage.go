package format

import (
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
)

// objectStorageOperationalStateColour maps managed object storage operational states to colours
func objectStorageOperationalStateColour(state upcloud.ManagedObjectStorageOperationalState) text.Colors {
	switch state {
	case upcloud.ManagedObjectStorageOperationalStateRunning:
		return text.Colors{text.FgGreen}
	case upcloud.ManagedObjectStorageOperationalStateDeleteDNS,
		upcloud.ManagedObjectStorageOperationalStateDeleteNetwork,
		upcloud.ManagedObjectStorageOperationalStatePending,
		upcloud.ManagedObjectStorageOperationalStateSetupCheckup,
		upcloud.ManagedObjectStorageOperationalStateSetupDNS,
		upcloud.ManagedObjectStorageOperationalStateSetupNetwork,
		upcloud.ManagedObjectStorageOperationalStateSetupService:
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// ObjectStorageOperationalState implements Format function for managed object storage operational states
func ObjectStorageOperationalState(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(objectStorageOperationalStateColour, val)
}

// objectStorageConfiguredStatusColour maps managed object storage configured statuses to colours
func objectStorageConfiguredStatusColour(state upcloud.ManagedObjectStorageConfiguredStatus) text.Colors {
	switch state {
	case upcloud.ManagedObjectStorageConfiguredStatusStarted:
		return text.Colors{text.FgGreen}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// ObjectStorageConfiguredStatus implements Format function for managed object storage configured statuses
func ObjectStorageConfiguredStatus(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(objectStorageConfiguredStatusColour, val)
}
