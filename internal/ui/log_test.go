package ui

import (
	"bytes"
	"fmt"
	"testing"

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

	assert.Contains(t, b.String(), msg)
	assert.Contains(t, b.String(), ": Done")
}
