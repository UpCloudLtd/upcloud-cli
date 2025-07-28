package resolver

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
)

const (
	helpUUIDName  = "<UUID/Name...>"
	helpUUIDTitle = "<UUID/Title...>"
)

// Resolver represents the most basic argument resolver, a function that accepts and argument and returns the resolved value(s).
type Resolver func(arg string) (resolved Resolved)

// ResolutionProvider is an interface for commands that provide resolution, either custom or the built-in ones
type ResolutionProvider interface {
	Get(ctx context.Context, svc service.AllServices) (Resolver, error)
	PositionalArgumentHelp() string
}

type CachingResolutionProvider[T any] interface {
	ResolutionProvider
	GetCached(uuid string) (T, error)
}

type MatchType int

const (
	MatchTypeExact           MatchType = 4
	MatchTypeCaseInsensitive MatchType = 3
	MatchTypeGlobPattern     MatchType = 2
	MatchTypePrefix          MatchType = 1
	MatchTypeNone            MatchType = 0
)

type Resolved struct {
	Arg     string
	matches map[string]MatchType
}

// AddMatch adds a match to the resolved value. If the match is already present, the highest match type is kept. I.e., exact match is kept over case insensitive match.
func (r *Resolved) AddMatch(uuid string, matchType MatchType) {
	if r.matches == nil {
		r.matches = make(map[string]MatchType)
	}

	current := r.matches[uuid]
	r.matches[uuid] = max(current, matchType)
}

func (r *Resolved) getAll() ([]string, MatchType) {
	var all []string
	for _, matchType := range []MatchType{
		MatchTypeExact,
		MatchTypeCaseInsensitive,
		MatchTypeGlobPattern,
		MatchTypePrefix,
	} {
		for uuid, match := range r.matches {
			if match == matchType {
				all = append(all, uuid)
			}
		}

		if len(all) > 0 {
			return all, matchType
		}
	}

	return all, MatchTypeNone
}

// GetAll returns matches with match-type that equals the highest available match-type for the resolved value. I.e., if there is an exact match, only exact matches are returned even if there would be case-insensitive matches.
//
// If match-type is not a glob pattern match, an error is returned if there are multiple matches.
func (r *Resolved) GetAll() ([]string, error) {
	all, matchType := r.getAll()

	if len(all) == 0 {
		return nil, NotFoundError(r.Arg)
	}

	// For backwards compatibility, allow multiple matches only for glob patterns.
	if len(all) > 1 && matchType != MatchTypeGlobPattern {
		return nil, NonGlobMultipleMatchesError(r.Arg)
	}
	return all, nil
}

// GetOnly returns the only match if there is only one match. If there are no or multiple matches, an empty value and an error is returned.
func (r *Resolved) GetOnly() (string, error) {
	all, _ := r.getAll()
	if len(all) == 0 {
		return "", NotFoundError(r.Arg)
	}

	if len(all) > 1 {
		return "", AmbiguousResolutionError(r.Arg)
	}

	return all[0], nil
}
