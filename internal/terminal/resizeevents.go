// +build linux darwin

package terminal

import (
	"os"
	"os/signal"
	"syscall"
)

// ResizeListener calls given callback function when terminal window is resized.
// This is currently a no-op in Windows.
type ResizeListener struct {
	signal chan os.Signal
}

// NewResizeListener creates a new ResizeListener
func NewResizeListener(callback func()) *ResizeListener {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGWINCH)
	go func() {
		for range signalCh {
			callback()
		}
	}()
	return &ResizeListener{
		signal: signalCh,
	}
}

// Close stops ResizeListener listening and cleans up associated resources
func (s *ResizeListener) Close() {
	signal.Stop(s.signal)
	close(s.signal)
}
