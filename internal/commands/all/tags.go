package all

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// cachingTag implements resolver for tags by name, caching the results
type cachingTag struct {
	resolver.Cache[upcloud.Tag]

	Type string
}

var (
	_ resolver.ResolutionProvider                     = &cachingTag{}
	_ resolver.CachingResolutionProvider[upcloud.Tag] = &cachingTag{}
)

// Get implements ResolutionProvider.Get
func (s *cachingTag) Get(ctx context.Context, svc internal.AllServices) (resolver.Resolver, error) {
	tags, err := svc.GetTags(ctx)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags.Tags {
		s.AddCached(tag.Name, tag)
	}

	return func(arg string) resolver.Resolved {
		rv := resolver.Resolved{Arg: arg}
		for _, tag := range tags.Tags {
			rv.AddMatch(tag.Name, resolver.MatchTitle(arg, tag.Name))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s *cachingTag) PositionalArgumentHelp() string {
	return "<Name...>"
}

func deleteTag(exec commands.Executor, name string) (output.Output, error) {
	svc := exec.All()
	err := svc.DeleteTag(exec.Context(), &request.DeleteTagRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return output.None{}, nil
}
