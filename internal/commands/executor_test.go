package commands

import (
	"fmt"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"

	"github.com/stretchr/testify/assert"
)

func TestExecutor_WaitFor(t *testing.T) {
	mService := &smock.Service{}
	cfg := config.New()
	exec := NewExecutor(cfg, mService, cfg.NewLogger("test"))
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
	mService := &smock.Service{}
	cfg := config.New()
	exec := NewExecutor(cfg, mService, cfg.NewLogger("test"))
	err := exec.WaitFor(func() error {
		time.Sleep(50 * time.Minute)
		return nil
	}, 100*time.Millisecond)
	assert.EqualError(t, err, "timed out")
}

func TestExecutor_WaitForError(t *testing.T) {
	mService := &smock.Service{}
	cfg := config.New()
	exec := NewExecutor(cfg, mService, cfg.NewLogger("test"))
	err := exec.WaitFor(func() error {
		time.Sleep(10 * time.Millisecond)
		return fmt.Errorf("mockmock")
	}, 100*time.Millisecond)
	assert.EqualError(t, err, "mockmock")
}
