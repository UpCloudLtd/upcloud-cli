package commands

import (
	"fmt"
	"testing"
	"time"

	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"

	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
)

func TestExecutor_WaitFor(t *testing.T) {
	mService := &smock.Service{}
	cfg := config.New()
	exec := NewExecutor(cfg, mService, flume.New("test"))
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
	exec := NewExecutor(cfg, mService, flume.New("test"))
	err := exec.WaitFor(func() error {
		time.Sleep(50 * time.Minute)
		return nil
	}, 100*time.Millisecond)
	assert.EqualError(t, err, "timed out")
}

func TestExecutor_WaitForError(t *testing.T) {
	mService := &smock.Service{}
	cfg := config.New()
	exec := NewExecutor(cfg, mService, flume.New("test"))
	err := exec.WaitFor(func() error {
		time.Sleep(10 * time.Millisecond)
		return fmt.Errorf("mockmock")
	}, 100*time.Millisecond)
	assert.EqualError(t, err, "mockmock")
}

type mockLogEntry struct {
	Msg  string
	Args []interface{}
}

type mockLogger struct {
	debugLines []mockLogEntry
	infoLines  []mockLogEntry
	errorLines []mockLogEntry
	context    []interface{}
}

func (m *mockLogger) Debug(msg string, args ...interface{}) {
	m.debugLines = append(m.debugLines, mockLogEntry{msg, args})
}

func (m *mockLogger) Info(msg string, args ...interface{}) {
	m.infoLines = append(m.infoLines, mockLogEntry{msg, args})
}

func (m *mockLogger) Error(msg string, args ...interface{}) {
	m.errorLines = append(m.errorLines, mockLogEntry{msg, args})
}

func (m mockLogger) IsDebug() bool {
	return true
}

func (m mockLogger) IsInfo() bool {
	return true
}

func (m mockLogger) With(args ...interface{}) flume.Logger {
	return &mockLogger{
		context: append(m.context, args...),
	}
}

func TestExecutor_Logging(t *testing.T) {
	mService := &smock.Service{}
	cfg := config.New()
	logger := &mockLogger{context: []interface{}{"base", "context"}}
	exec := NewExecutor(cfg, mService, logger)
	exec.LogInfo("info", "hello", "world")
	exec.LogDebug("debug", "hi", "earth")
	assert.Equal(t, &mockLogger{
		debugLines: []mockLogEntry{
			{Msg: "debug", Args: []interface{}{"hi", "earth"}},
		},
		infoLines: []mockLogEntry{
			{Msg: "info", Args: []interface{}{"hello", "world"}},
		},
		context: []interface{}{"base", "context"},
	}, logger)
	// test adding context
	contextExec := exec.WithLogger("added", "newcontext")
	contextExec.LogInfo("infoz", "helloz", "worldz")
	assert.Equal(t, &mockLogger{
		infoLines: []mockLogEntry{
			{Msg: "infoz", Args: []interface{}{"helloz", "worldz"}},
		},
		context: []interface{}{"base", "context", "added", "newcontext"},
	}, contextExec.(*executorImpl).logger)
}
