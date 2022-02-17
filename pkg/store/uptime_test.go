package store_test

import (
	"encoding/json"
	"fmt"
	libKVStore "github.com/abronan/valkeyrie/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockKVStore "gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/internal/mocks/valkeyrie/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/testutils"
	"testing"
)

func TestPersistentStore_Uptime(t *testing.T) {
	t.Run("should get uptime", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		now := testutils.Uptime()
		data, _ := json.Marshal(now)
		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/status/server", nil).Return(&libKVStore.KVPair{Value: data}, nil)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		since, err := persistentStore.Uptime()

		assert.NoError(t, err)
		assert.Equal(t, now, since)
	})

	t.Run("should return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/status/server", nil).Return(nil, fmt.Errorf("store not available"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		since, err := persistentStore.Uptime()

		if assert.Error(t, err) {
			assert.Equal(t, "store not available", err.Error())
		}
		assert.Nil(t, since)
	})
}

func TestPersistentStore_RefreshUptime(t *testing.T) {
	t.Run("should refresh uptime", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Put("dwarka/status/server", gomock.Any(), nil).Return(nil)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.RefreshUptime()

		assert.NoError(t, err)
	})

	t.Run("should return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Put("dwarka/status/server", gomock.Any(), nil).Return(fmt.Errorf("unable to save"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.RefreshUptime()

		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})
}
