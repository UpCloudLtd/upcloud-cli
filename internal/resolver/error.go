package resolver

import "fmt"

// AmbiguousResolutionError is a resolver error when multiple matching entries have been found.
type AmbiguousResolutionError string

var _ error = AmbiguousResolutionError("")

func (s AmbiguousResolutionError) Error() string {
	return fmt.Sprintf("'%v' is ambiguous, found multiple matches", string(s))
}

// NonGlobMultipleMatchesError is a resolver error when multiple matching entries have been found with non-glob argument.
type NonGlobMultipleMatchesError string

var _ error = NonGlobMultipleMatchesError("")

func (s NonGlobMultipleMatchesError) Error() string {
	return fmt.Sprintf("'%v' is not a glob pattern, but matches multiple values. To target multiple resources with single argument, use a glob pattern, e.g. server-*", string(s))
}

// NotFoundError is a resolver error when no matching entries have been found.
type NotFoundError string

var _ error = NotFoundError("")

func (s NotFoundError) Error() string {
	return fmt.Sprintf("nothing found matching '%v'", string(s))
}

type CacheUninitializedError string

var _ error = CacheUninitializedError("")

func (s CacheUninitializedError) Error() string {
	return fmt.Sprintf("resolver cache for %s has not been initialized", string(s))
}
