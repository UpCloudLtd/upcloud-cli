package resolver

import "fmt"

// AmbiguousResolutionError is a resolver error when multiple matching entries have been found
type AmbiguousResolutionError string

var _ error = AmbiguousResolutionError("")

func (s AmbiguousResolutionError) Error() string {
	return fmt.Sprintf("'%v' is ambiguous, found multiple matches", string(s))
}

// NotFoundError is a resolver error when no matching entries have been found
type NotFoundError string

var _ error = NotFoundError("")

func (s NotFoundError) Error() string {
	return fmt.Sprintf("nothing found matching '%v'", string(s))
}
