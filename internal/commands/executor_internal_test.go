package commands

import (
	"testing"

	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"

	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
)

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
	t.Parallel()
	mService := &smock.Service{}
	cfg := config.New()
	logger := &mockLogger{context: []interface{}{"base", "context"}}
	exec := NewExecutor(cfg, mService, logger)
	exec.Debug("debug1", "hello", "world")
	// create a contexted executor
	contextExec := exec.WithLogger("added", "newcontext")
	contextExec.Debug("debugcontext", "helloz", "worldz")
	exec.Debug("debug2", "hi", "earth")
	// make sure the main executor does not leak to the contexted one or vice versa
	assert.Equal(t, &mockLogger{
		debugLines: []mockLogEntry{
			{Msg: "debug1", Args: []interface{}{"hello", "world"}},
			{Msg: "debug2", Args: []interface{}{"hi", "earth"}},
		},
		context: []interface{}{"base", "context"},
	}, logger)
	assert.Equal(t, &mockLogger{
		debugLines: []mockLogEntry{
			{Msg: "debugcontext", Args: []interface{}{"helloz", "worldz"}},
		},
		context: []interface{}{"base", "context", "added", "newcontext"},
	}, contextExec.(*executorImpl).logger)
}
