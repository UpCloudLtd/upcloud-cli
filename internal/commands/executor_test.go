package commands_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"

	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
)

func TestExecutor_WaitFor(t *testing.T) {
	t.Parallel()
	mService := &smock.Service{}
	cfg := config.New()
	exec := commands.NewExecutor(cfg, mService, flume.New("test"))
	finished := false
	// normal operation
	go func() {
		err := exec.WaitFor(func() error {
			time.Sleep(50 * time.Millisecond)
			finished = true
			return nil
		}, time.Minute)
		assert.NoError(t, err)
	}()
	assert.Eventually(t, func() bool {
		return finished
	}, time.Second, 10*time.Millisecond)
}

func TestExecutor_WaitForTimeout(t *testing.T) {
	t.Parallel()
	mService := &smock.Service{}
	cfg := config.New()
	exec := commands.NewExecutor(cfg, mService, flume.New("test"))
	err := exec.WaitFor(func() error {
		time.Sleep(50 * time.Minute)
		return nil
	}, 100*time.Millisecond)
	assert.EqualError(t, err, "timed out")
}

func TestExecutor_WaitForError(t *testing.T) {
	t.Parallel()
	mService := &smock.Service{}
	cfg := config.New()
	exec := commands.NewExecutor(cfg, mService, flume.New("test"))
	err := exec.WaitFor(func() error {
		time.Sleep(10 * time.Millisecond)
		return fmt.Errorf("mockmock")
	}, 100*time.Millisecond)
	assert.EqualError(t, err, "mockmock")
}
