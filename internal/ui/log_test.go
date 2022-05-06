package ui

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

// TestLogRender ensures that log entries are outputted with single Render call
func TestLogRender(t *testing.T) {
	var b bytes.Buffer
	config := LiveLogDefaultConfig
	config.EntryMaxWidth = 50
	livelog := NewLiveLog(&b, config)

	// Simulate typical livelog usage
	msg := "Testing log.Render"
	entry := NewLogEntry(msg)
	livelog.AddEntries(entry)
	entry.StartedNow()
	entry.SetMessage(fmt.Sprintf("%s: Done", msg))
	entry.MarkDone()
	livelog.Render()
	livelog.Close()

	assert.Equal(t, entry.result, logEntrySuccess)
	assert.Contains(t, b.String(), msg)
	assert.Contains(t, b.String(), ": Done")
}

func TestLogErrors(t *testing.T) {
	text.EnableColors()

	var b bytes.Buffer
	config := LiveLogDefaultConfig
	config.EntryMaxWidth = 50
	livelog := NewLiveLog(&b, config)

	failedEntry := NewLogEntry("MarkFailed")
	warningEntry := NewLogEntry("MarkWarning")
	livelog.AddEntries(failedEntry, warningEntry)

	failedEntry.StartedNow()
	warningEntry.StartedNow()
	livelog.Render()

	failedEntry.MarkFailed()
	warningEntry.MarkWarning()
	livelog.Render()
	livelog.Close()

	assert.Equal(t, failedEntry.result, logEntryFailed)
	assert.Equal(t, warningEntry.result, logEntryWarning)

	// Test that output contains partial ANSI code before message
	assert.Contains(t, b.String(), `[91mMarkFailed`)
	assert.Contains(t, b.String(), `[93mMarkWarning`)
}
