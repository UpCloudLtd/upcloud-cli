package namedargs

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// GetStorage finds storage based on values provided to named args (e.g., --storage storage-name)
func GetStorage(exec commands.Executor, arg string) (upcloud.Storage, error) {
	return Get(&resolver.CachingStorage{}, exec, arg)
}
