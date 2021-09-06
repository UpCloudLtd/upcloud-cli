package resolver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
)

func testResolutionFailure(t *testing.T, methodName string, response interface{}, testResolver resolver.ResolutionProvider, ambiguousCases []string, notFoundCases []string) {
	t.Helper()
	mService := &smock.Service{}
	mService.On(methodName).Return(response, nil)
	argResolver, err := testResolver.Get(mService)
	assert.NoError(t, err)
	for _, ambiguousCase := range ambiguousCases {
		resolved, err := argResolver(ambiguousCase)
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(ambiguousCase))
		assert.Equal(t, "", resolved)
	}
	for _, notFoundCase := range notFoundCases {
		resolved, err := argResolver(notFoundCase)
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.NotFoundError(notFoundCase))
		assert.Equal(t, "", resolved)
	}
	// make sure caching works, eg. we didn't call the method to be cached more than once
	mService.AssertNumberOfCalls(t, methodName, 1)
}
