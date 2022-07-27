package clierrors

const UnspecifiedErrorCode int = 100

// ClientError declares interface for errors known to the client that set specific error code.
type ClientError interface {
	ErrorCode() int
}
