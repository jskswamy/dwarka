package api_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockStore "gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/internal/mocks/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/testutils"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	t.Run("should return server start time with http.StatusOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockKVStore := mockStore.NewMockStore(ctrl)
		mockKVStore.EXPECT().Uptime().Return(testutils.Uptime(), nil)

		request, err := http.NewRequest("GET", "http://test/ping", nil)
		if err != nil {
			t.Error(err)
		}

		res, err := testutils.ServeHTTPRequest(mockKVStore, request)
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("should handle error from store with http.StatusInternalServerError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockKVStore := mockStore.NewMockStore(ctrl)
		mockKVStore.EXPECT().Uptime().Return(nil, fmt.Errorf("unable to read value from store"))

		request, err := http.NewRequest("GET", "http://test/ping", nil)
		if err != nil {
			t.Error(err)
		}

		res, err := testutils.ServeHTTPRequest(mockKVStore, request)
		assert.NoError(t, err)
		assert.Equal(t, 500, res.StatusCode)

		expected := "there was problem when reading value from store, reason: unable to read value from store"
		msg, err := testutils.ReadError(res)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, msg)
		}
	})
}
