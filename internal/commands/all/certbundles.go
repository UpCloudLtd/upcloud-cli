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

// cachingCertificateBundle implements resolver for certificate bundles by name, caching the results
type cachingCertificateBundle struct {
	resolver.Cache[upcloud.LoadBalancerCertificateBundle]

	Type string
}

var (
	_ resolver.ResolutionProvider                                               = &cachingCertificateBundle{}
	_ resolver.CachingResolutionProvider[upcloud.LoadBalancerCertificateBundle] = &cachingCertificateBundle{}
)

// Get implements ResolutionProvider.Get
func (s *cachingCertificateBundle) Get(ctx context.Context, svc internal.AllServices) (resolver.Resolver, error) {
	certBundles, err := svc.GetLoadBalancerCertificateBundles(ctx, &request.GetLoadBalancerCertificateBundlesRequest{
		Page: &request.Page{
			Number: 0,
			Size:   100,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, cb := range certBundles {
		s.AddCached(cb.UUID, cb)
	}

	return func(arg string) resolver.Resolved {
		rv := resolver.Resolved{Arg: arg}
		for _, cb := range certBundles {
			rv.AddMatch(cb.UUID, resolver.MatchTitle(arg, cb.Name))
			rv.AddMatch(cb.UUID, resolver.MatchUUID(arg, cb.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s *cachingCertificateBundle) PositionalArgumentHelp() string {
	return "<Name...>"
}

func deleteCertificateBundle(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	err := svc.DeleteLoadBalancerCertificateBundle(exec.Context(), &request.DeleteLoadBalancerCertificateBundleRequest{
		UUID: uuid,
	})
	if err != nil {
		return nil, err
	}

	return output.None{}, nil
}
