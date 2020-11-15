package main

import (
	"github.com/UpCloudLtd/cli/internal/commands/storage"
	"testing"
)

func Test(t *testing.T) {
	_ = mc.Cobra().RunE(storage.StorageCommand().Cobra(), []string{"storage", "list"})
}
