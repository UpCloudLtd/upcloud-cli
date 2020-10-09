package ui

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/cli/internal/terminal"
)

var (
	LiveLogDefaultColours = LiveLogColours{
		Pending:    text.Colors{text.FgHiBlack},
		InProgress: text.Colors{text.FgHiBlue, text.Bold},
		Done:       text.Colors{text.FgHiGreen},
		Details:    text.Colors{text.FgHiBlack},
		Time:       text.Colors{text.FgHiCyan},
	}
	LiveLogDefaultConfig = LiveLogConfig{
		EntryMaxWidth:        80,
		RenderPending:        true,
		DisableLiveRendering: !terminal.IsStdoutTerminal(),
		Colours:              LiveLogDefaultColours,
	}
	LiveLogEntryErrorColours = text.FgHiRed
)

type LiveLogColours struct {
	Pending    text.Colors
	InProgress text.Colors
	Done       text.Colors
	Details    text.Colors
	Time       text.Colors
}

type LiveLogConfig struct {
	EntryMaxWidth        int
	RenderPending        bool
	DisableLiveRendering bool
	Colours              LiveLogColours
}

func NewLiveLog(out io.Writer, style LiveLogConfig) *LiveLog {
	return &LiveLog{out: out, config: style}
}

type LiveLog struct {
	mu                sync.Mutex
	config            LiveLogConfig
	entriesPending    []*LogEntry
	entriesInProgress []*LogEntry
	entriesDone       []*LogEntry
	renderingStarted  bool
	height            int
	out               io.Writer
	isTerminal        bool
}

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

func (s *LiveLog) Render() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Acquire locks for all entries
	locks := s.lockActiveEntries()
	defer s.unlockActiveEntries(locks)

	if s.height > 0 {
		// Set cursor to start of the current rendering space
		s.write(text.CursorUp.Sprintn(s.height))
	}
	// fmt.Println(len(s.entriesPending), len(s.entriesInProgress), len(s.entriesDone))
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
			details := text.WrapSoft(entry.details, s.config.EntryMaxWidth-len(entry.detailsPrefix))
			details = text.Pad(details, s.config.EntryMaxWidth, ' ')
			s.write(s.config.Colours.Details.Sprint(IndentText(details, entry.detailsPrefix, true)))
			s.write("\n")
		}
		s.eraseLine()
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

	// Render in-progress entries
	s.height = len(s.entriesInProgress)
	for _, entry := range s.entriesInProgress {
		if entry.done {
			continue
		}
		s.renderEntry(entry)
	}

	// Render queued
	if s.config.RenderPending {
		s.height += len(s.entriesPending)
		for _, entry := range s.entriesPending {
			if entry.started.IsZero() {
				s.renderEntry(entry)
			}
		}
	}
}

func (s *LiveLog) renderEntry(entry *LogEntry) {
	s.eraseLine()
	var durStr string
	var colours text.Colors
	switch {
	case entry.started.IsZero():
		colours = s.config.Colours.Pending
	case entry.done:
		colours = s.config.Colours.Done
	default:
		colours = s.config.Colours.InProgress
	}
	if !entry.started.IsZero() {
		dur := time.Now().Sub(entry.started)
		durStr = fmt.Sprintf("%dm%ds", int(dur.Minutes()), int(dur.Seconds()))
	}
	msg := entry.msg
	if text.RuneCount(msg) > s.config.EntryMaxWidth {
		msg = fmt.Sprintf("%s...", text.Trim(msg, s.config.EntryMaxWidth-3))
	}
	msg = text.Pad(msg, s.config.EntryMaxWidth, ' ')
	s.write(colours.Sprint(msg))
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

func NewLogEntry(msg string) *LogEntry {
	return &LogEntry{msg: msg}
}

type LogEntry struct {
	mu            sync.Mutex
	msg           string
	details       string
	detailsPrefix string
	started       time.Time
	done          bool
}

func (s *LogEntry) SetMessage(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.msg = msg
}

func (s *LogEntry) SetDetails(details, prefix string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.details = details
	s.detailsPrefix = prefix
}

func (s *LogEntry) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.started = time.Now()
}

func (s *LogEntry) MarkDone() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.done = true
}
