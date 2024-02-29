package resolver_test

import (
	"context"
	"errors"
	"testing"

	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockGateways = []upcloud.Gateway{
	{Name: "asd", UUID: "abcdef"},
	{Name: "asd", UUID: "abcghi"},
	{Name: "qwe", UUID: "jklmno"},
}

func TestGatewayResolution(t *testing.T) {
	t.Run("UUID", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetGateways", mock.Anything).Return(mockGateways, nil)
		res := resolver.CachingGateway{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, db := range mockGateways {
			resolved, err := argResolver(db.UUID)
			assert.NoError(t, err)
			assert.Equal(t, db.UUID, resolved)
		}

		// Make sure caching works, eg. we didn't call GetGateways more than once
		mService.AssertNumberOfCalls(t, "GetGateways", 1)
	})

	t.Run("Name", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetGateways", mock.Anything).Return(mockGateways, nil)
		res := resolver.CachingGateway{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		db := mockGateways[2]
		resolved, err := argResolver(db.Name)
		assert.NoError(t, err)
		assert.Equal(t, db.UUID, resolved)
		// Make sure caching works, eg. we didn't call GetGateways more than once
		mService.AssertNumberOfCalls(t, "GetGateways", 1)
	})

	t.Run("Failures", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetGateways", mock.Anything).Return(mockGateways, nil)

		res := resolver.CachingGateway{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		var resolved string

		// Ambigous Name
		resolved, err = argResolver("asd")
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError("asd"))
		assert.Equal(t, "", resolved)

		// Not found
		resolved, err = argResolver("not-found")
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.NotFoundError("not-found"))
		assert.Equal(t, "", resolved)

		// Make sure caching works, eg. we didn't call GetGateways more than once
		mService.AssertNumberOfCalls(t, "GetGateways", 1)
	})
}

func TestFailingGatewayResolution(t *testing.T) {
	mService := &smock.Service{}
	mService.On("GetGateways", mock.Anything).Return(nil, errors.New("MOCKERROR"))
	res := resolver.CachingGateway{}
	argResolver, err := res.Get(context.TODO(), mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}
