package utils

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

func propertyKey(keySlice ...string) string {
	str := keySlice[0]
	for _, i := range keySlice[1:] {
		str = fmt.Sprintf("%s.%s", str, i)
	}
	return str
}

func GetFlatDatabaseProperties(properties map[string]upcloud.ManagedDatabaseServiceProperty, keyPrefix ...string) map[string]upcloud.ManagedDatabaseServiceProperty {
	flat := make(map[string]upcloud.ManagedDatabaseServiceProperty)
	for key, details := range properties {
		keySlice := append(keyPrefix, key) //nolint:gocritic // Construct key slice from prefix set by parent and key from current property
		flat[propertyKey(keySlice...)] = details
		for k, v := range GetFlatDatabaseProperties(details.Properties, keySlice...) {
			flat[k] = v
		}
	}

	return flat
}
