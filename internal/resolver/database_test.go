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

var mockDatabases = []upcloud.ManagedDatabase{
	{Title: "asd", UUID: "abcdef"},
	{Title: "asd", UUID: "abcghi"},
	{Title: "qwe", UUID: "jklmno"},
}

func TestDatabaseResolution(t *testing.T) {
	t.Run("UUID", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetManagedDatabases", mock.Anything).Return(mockDatabases, nil)
		res := resolver.CachingDatabase{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, db := range mockDatabases {
			resolved := argResolver(db.UUID)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, db.UUID, value)
		}

		// Make sure caching works, eg. we didn't call GetManagedDatabases more than once
		mService.AssertNumberOfCalls(t, "GetManagedDatabases", 1)
	})

	t.Run("Title", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetManagedDatabases", mock.Anything).Return(mockDatabases, nil)
		res := resolver.CachingDatabase{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		db := mockDatabases[2]
		resolved := argResolver(db.Title)
		value, err := resolved.GetOnly()
		assert.NoError(t, err)
		assert.Equal(t, db.UUID, value)
		// Make sure caching works, eg. we didn't call GetManagedDatabases more than once
		mService.AssertNumberOfCalls(t, "GetManagedDatabases", 1)
	})

	t.Run("Failures", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetManagedDatabases", mock.Anything).Return(mockDatabases, nil)

		res := resolver.CachingDatabase{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		// Ambiguous title
		resolved := argResolver("asd")
		value, err := resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError("asd"))
		assert.Equal(t, "", value)

		// Not found
		resolved = argResolver("not-found")
		value, err = resolved.GetOnly()

		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.NotFoundError("not-found"))
		assert.Equal(t, "", value)

		// Make sure caching works, eg. we didn't call GetManagedDatabases more than once
		mService.AssertNumberOfCalls(t, "GetManagedDatabases", 1)
	})
}

func TestFailingDatabaseResolution(t *testing.T) {
	mService := &smock.Service{}
	mService.On("GetManagedDatabases", mock.Anything).Return(nil, errors.New("MOCKERROR"))
	res := resolver.CachingDatabase{}
	argResolver, err := res.Get(context.TODO(), mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}
