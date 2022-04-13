package resolver_test

import (
	"errors"
	"testing"

	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	s1   = upcloud.Storage{UUID: "abc", Title: "adads"}
	s2   = upcloud.Storage{UUID: "bcd", Title: "badada"}
	s3   = upcloud.Storage{UUID: "cde", Title: "cdasds"}
	amb1 = upcloud.Storage{UUID: "def", Title: "dadads"}
	amb2 = upcloud.Storage{UUID: "def", Title: "eadads"}
	amb3 = upcloud.Storage{UUID: "ghi", Title: "dadads"}
)
var mockStorages = &upcloud.Storages{Storages: []upcloud.Storage{s1, s2, s3, amb1, amb2, amb3}}
var unambiguousStorages = []upcloud.Storage{s1, s2, s3}

func TestStorageResolution(t *testing.T) {
	t.Run("resolve uuid", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetStorages", mock.Anything).Return(mockStorages, nil)
		res := resolver.CachingStorage{}
		argResolver, err := res.Get(mService)
		assert.NoError(t, err)
		for _, storage := range unambiguousStorages {
			resolved, err := argResolver(storage.UUID)
			assert.NoError(t, err)
			assert.Equal(t, storage.UUID, resolved)
		}
		// make sure caching works, eg. we didn't call GetStorages more than once
		mService.AssertNumberOfCalls(t, "GetStorages", 1)
	})

	t.Run("resolve title", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetStorages", mock.Anything).Return(mockStorages, nil)
		res := resolver.CachingStorage{}
		argResolver, err := res.Get(mService)
		assert.NoError(t, err)
		for _, storage := range unambiguousStorages {
			resolved, err := argResolver(storage.Title)
			assert.NoError(t, err)
			assert.Equal(t, storage.UUID, resolved)
		}
		// make sure caching works, eg. we didn't call GetStorages more than once
		mService.AssertNumberOfCalls(t, "GetStorages", 1)
	})

	t.Run("failure situations", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetStorages", mock.Anything).Return(mockStorages, nil)

		res := resolver.CachingStorage{}
		argResolver, err := res.Get(mService)
		assert.NoError(t, err)

		// ambigous uuid
		resolved, err := argResolver(amb2.UUID)
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(amb2.UUID))
		assert.Equal(t, "", resolved)

		// ambigous title
		resolved, err = argResolver(amb1.Title)
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(amb1.Title))
		assert.Equal(t, "", resolved)

		// not found
		resolved, err = argResolver("notfounf")
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.NotFoundError("notfounf"))
		assert.Equal(t, "", resolved)

		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetStorages", 1)
	})
}

func TestFailingStorageResolution(t *testing.T) {
	mService := &smock.Service{}
	mService.On("GetStorages", mock.Anything).Return(nil, errors.New("MOCKERROR"))
	res := resolver.CachingStorage{}
	argResolver, err := res.Get(mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}
