package ui

import (
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/upcloud-cli/internal/terminal"
)

var (
	liveLogDefaultColours = liveLogColours{
		Pending:    text.Colors{text.FgHiBlack},
		InProgress: text.Colors{text.FgHiBlue, text.Bold},
		Done:       text.Colors{text.FgHiGreen},
		Details:    text.Colors{text.FgHiBlack},
		Time:       text.Colors{text.FgHiCyan},
	}
	// LiveLogDefaultConfig represents the default settings for live log.
	LiveLogDefaultConfig = LiveLogConfig{
		EntryMaxWidth:        0,
		renderPending:        true,
		DisableLiveRendering: !terminal.IsStdoutTerminal(),
		Colours:              liveLogDefaultColours,
	}
	// LiveLogEntryErrorColours specifies the colour used for errors in LiveLog
	// TODO: remove cross-package dependency and make this private.
	LiveLogEntryErrorColours = text.FgHiRed
)

type liveLogColours struct {
	Pending    text.Colors
	InProgress text.Colors
	Done       text.Colors
	Details    text.Colors
	Time       text.Colors
}

// LiveLogConfig is a configuration for rendering live logs.
type LiveLogConfig struct {
	// EntryMaxWidth sets the forced maximum width of live log lines. Use 0 for detecting terminal width.
	EntryMaxWidth        int
	renderPending        bool
	DisableLiveRendering bool
	Colours              liveLogColours
}

// LiveLog represents the internal state of a live log renderer.
type LiveLog struct {
	mu                sync.Mutex
	config            LiveLogConfig
	entriesPending    []*LogEntry
	entriesInProgress []*LogEntry
	entriesDone       []*LogEntry
	lastRenderWidth   int
	height            int
	out               io.Writer
	resizeListener    *terminal.ResizeListener
}

// NewLiveLog returns a new renderer for live logs.
func NewLiveLog(out io.Writer, style LiveLogConfig) *LiveLog {
	llog := &LiveLog{out: out, config: style}
	llog.resizeListener = terminal.NewResizeListener(llog.Render)
	return llog
}

// Close closes the LiveLog and cleans up related resources.
func (s *LiveLog) Close() {
	s.resizeListener.Close()
}

// AddEntries adds log entries to LiveLog.
func (s *LiveLog) AddEntries(entries ...*LogEntry) {
	for _, e := range entries {
		if e == nil {
			panic("LiveLog: nil entry given")
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entriesPending = append(s.entriesPending, entries...)
}

func (s *LiveLog) lockActiveEntries() []*sync.Mutex {
	var locks []*sync.Mutex
	for _, entry := range s.entriesPending {
		entry.mu.Lock()
		locks = append(locks, &entry.mu)
	}
	for _, entry := range s.entriesInProgress {
		entry.mu.Lock()
		locks = append(locks, &entry.mu)
	}
	return locks
}

func (s *LiveLog) unlockActiveEntries(locks []*sync.Mutex) {
	for _, mu := range locks {
		mu.Unlock()
	}
}

// Render renders the LiveLog.
func (s *LiveLog) Render() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Acquire locks for all entries
	locks := s.lockActiveEntries()
	defer s.unlockActiveEntries(locks)

	// Set cursor to start of the current rendering space
	heightDelta := 0
	if s.lastRenderWidth > 0 && s.lastRenderWidth != s.MaxWidth() {
		// terminal size has changed since our last render, figure out where to actually move to
		prevBuffLen := s.height * s.lastRenderWidth
		expectedHeight := int(math.Ceil(float64(prevBuffLen) / float64(s.MaxWidth())))
		heightDelta = expectedHeight - s.height
	}
	// linefeed to move to the first column
	s.write("\r")
	// and clear up our rendering space
	for n := 0; n < s.height+heightDelta; n++ {
		s.eraseLine()                   // erase the end of line
		s.write(text.CursorUp.Sprint()) // and move up
	}
	newInProgress := s.entriesInProgress[:0]
	// Render any completed
	for _, entry := range s.entriesInProgress {
		if !entry.done {
			newInProgress = append(newInProgress, entry)
			continue
		}
		s.entriesDone = append(s.entriesDone, entry)
		s.renderEntry(entry)
		if entry.details != "" {
			details := text.WrapSoft(entry.details, s.MaxWidth()-len(entry.detailsPrefix))
			details = text.Pad(details, s.MaxWidth(), ' ')
			s.write(s.config.Colours.Details.Sprint(IndentText(details, entry.detailsPrefix, true)))
			s.write("\n")
		}
	}
	s.entriesInProgress = newInProgress

	// Add any pending entries that have started
	newPending := s.entriesPending[:0]
	for _, entry := range s.entriesPending {
		isStarted := !entry.started.IsZero()
		if !isStarted {
			newPending = append(newPending, entry)
			continue
		}
		s.entriesInProgress = append(s.entriesInProgress, entry)
	}
	s.entriesPending = newPending

	if s.config.DisableLiveRendering {
		return
	}

	s.height = 0
	// Render in-progress entries
	for _, entry := range s.entriesInProgress {
		if entry.done {
			continue
		}
		s.height++
		s.renderEntry(entry)
	}

	// Render queued
	if s.config.renderPending {
		for _, entry := range s.entriesPending {
			if entry.started.IsZero() {
				s.height++
				s.renderEntry(entry)
			}
		}
	}
	s.lastRenderWidth = s.MaxWidth()
}

func (s *LiveLog) renderEntry(entry *LogEntry) {
	var durStr string
	var lineColour text.Colors
	switch {
	case entry.started.IsZero():
		lineColour = s.config.Colours.Pending
	case entry.done:
		lineColour = s.config.Colours.Done
	default:
		lineColour = s.config.Colours.InProgress
	}
	if !entry.started.IsZero() {
		dur := time.Since(entry.started)
		durStr = fmt.Sprintf("%dm%ds",
			dur/time.Minute,
			dur%time.Minute/time.Second)
	}
	msg := entry.msg
	maxWidth := s.MaxWidth() - 1 - text.RuneCount(durStr)
	if text.RuneCount(msg) > maxWidth {
		msg = fmt.Sprintf("%s...", text.Trim(msg, maxWidth-3))
	}
	msg = text.Pad(msg, maxWidth, ' ')
	s.write(lineColour.Sprint(msg))
	s.write(s.config.Colours.Time.Sprint(durStr))
	s.write("\n")
}

func (s *LiveLog) eraseLine() {
	if s.config.DisableLiveRendering {
		return
	}
	s.write(text.EraseLine.Sprint())
}

func (s *LiveLog) write(str string) {
	_, err := fmt.Fprint(s.out, str)
	if err != nil {
		panic(fmt.Sprintf("LiveLog rendering error: %v", err))
	}
}

// MaxWidth returns the maximum allowed width of a log line for this LiveLog instance.
// if config.MaxEntryWidth is set to 0, it will return the maximum terminal width.
// Note that this means that MaxWidth() can return different values across Render()s,
// as terminal could be resized between them.
func (s *LiveLog) MaxWidth() int {
	if s.config.EntryMaxWidth != 0 {
		return s.config.EntryMaxWidth
	}
	return terminal.GetTerminalWidth()
}

// NewLogEntry creates a new LogEntry with the specified message.
func NewLogEntry(msg string) *LogEntry {
	return &LogEntry{msg: msg}
}

// LogEntry represents a single log entry in a live log.
type LogEntry struct {
	mu            sync.Mutex
	msg           string
	details       string
	detailsPrefix string
	started       time.Time
	done          bool
}

// SetMessage sets the message of the LogEntry.
func (s *LogEntry) SetMessage(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.msg = msg
}

// SetDetails sets the details of the LogEntry.
func (s *LogEntry) SetDetails(details, prefix string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.details = details
	s.detailsPrefix = prefix
}

// StartedNow sets the LogEntry time to time.Now.
func (s *LogEntry) StartedNow() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.started = time.Now()
}

// MarkDone marks the LogEntry as done.
func (s *LogEntry) MarkDone() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.done = true
}
