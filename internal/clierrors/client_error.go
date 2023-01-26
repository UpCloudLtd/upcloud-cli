package clierrors

const (
	UnspecifiedErrorCode int = 100
	InterruptSignalCode  int = 101
)

// ClientError declares interface for errors known to the client that set specific error code.
type ClientError interface {
	ErrorCode() int
}
