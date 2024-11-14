package resolver

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
)

const helpUUIDTitle = "<UUID/Title...>"

// Resolver represents the most basic argument resolver, a function that accepts and argument and returns the resolved value(s).
type Resolver func(arg string) (resolved Resolved)

// ResolutionProvider is an interface for commands that provide resolution, either custom or the built-in ones
type ResolutionProvider interface {
	Get(ctx context.Context, svc service.AllServices) (Resolver, error)
	PositionalArgumentHelp() string
}

type MatchType int

const (
	MatchTypeExact           MatchType = 3
	MatchTypeCaseInsensitive MatchType = 2
	MatchTypeWildCard        MatchType = 2
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

// GetAll returns all matches with match-type that equals the highest available match-type for the resolved value. I.e., if there is an exact match, only exact matches are returned even if there would be case-insensitive matches.
func (r *Resolved) GetAll() ([]string, error) {
	var all []string
	for _, matchType := range []MatchType{
		MatchTypeExact,
		MatchTypeCaseInsensitive,
		MatchTypeWildCard,
		MatchTypePrefix,
	} {
		for uuid, match := range r.matches {
			if match == matchType {
				all = append(all, uuid)
			}
		}

		if len(all) > 0 {
			return all, nil
		}
	}

	var err error
	if len(all) == 0 {
		err = NotFoundError(r.Arg)
	}
	return all, err
}

// GetOnly returns the only match if there is only one match. If there are no or multiple matches, an empty value and an error is returned.
func (r *Resolved) GetOnly() (string, error) {
	all, err := r.GetAll()
	if err != nil {
		return "", err
	}

	if len(all) > 1 {
		return "", AmbiguousResolutionError(r.Arg)
	}

	return all[0], nil
}
