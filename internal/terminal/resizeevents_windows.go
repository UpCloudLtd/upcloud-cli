package terminal

// ResizeListener calls given callback function when terminal window is resized.
// This is currently a no-op in Windows.
type ResizeListener struct {
}

// NewResizeListener creates a new ResizeListener
func NewResizeListener(_ func()) *ResizeListener {
	return &ResizeListener{}
}

// Close stops ResizeListener listening and cleans up associated resources
func (s *ResizeListener) Close() {
}
