package resolver_test

import (
	"context"
	"errors"
	"testing"

	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestAccountResolution(t *testing.T) {
	Account1 := upcloud.AccountListItem{
		Username: "account-1",
	}

	Account2 := upcloud.AccountListItem{
		Username: "account-2",
	}

	Account3 := upcloud.AccountListItem{
		Username: "account-3",
	}

	Account4 := upcloud.AccountListItem{
		Username: "account-4",
	}

	allAccounts := upcloud.AccountList{
		Account1,
		Account2,
		Account3,
		Account4,
	}

	t.Run("resolve username", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetAccountList").Return(allAccounts, nil)
		res := resolver.CachingAccount{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, account := range allAccounts {
			resolved := argResolver(account.Username)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, account.Username, value)
		}
		// make sure caching works, eg. we didn't call GetAccountList more than once
		mService.AssertNumberOfCalls(t, "GetAccountList", 1)
	})

	t.Run("failure situations", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetAccountList").Return(allAccounts, nil)

		res := resolver.CachingAccount{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		// not found
		resolved := argResolver("notfound")
		value, err := resolved.GetOnly()

		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.NotFoundError("notfound"))
		assert.Equal(t, "", value)

		// make sure caching works, eg. we didn't call GetAccountList more than once
		mService.AssertNumberOfCalls(t, "GetAccountList", 1)
	})
}

func TestFailingAccountResolution(t *testing.T) {
	mService := &smock.Service{}
	mService.On("GetAccountList").Return(upcloud.AccountList{}, errors.New("MOCKERROR"))
	res := resolver.CachingAccount{}
	argResolver, err := res.Get(context.TODO(), mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}
