package clierrors

const (
	UnspecifiedErrorCode int = 100
	InterruptSignalCode  int = 101
	MissingCredentials   int = 102
	InvalidCredentials   int = 103
)

// ClientError declares interface for errors known to the client that set specific error code.
type ClientError interface {
	ErrorCode() int
}
