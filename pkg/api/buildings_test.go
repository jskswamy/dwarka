package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	mockStore "gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/internal/mocks/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/testutils"
	"net/http"
	"testing"
)

func TestBuildings(t *testing.T) {

	t.Run("test GET /buildings", func(t *testing.T) {

		building := gateway.Building{
			Lat:            1.2,
			Lan:            1.3,
			PhysicalEntity: gateway.PhysicalEntity{Name: "building one", Description: "test building"},
		}
		buildings := gateway.Buildings{building.ID(): building}

		t.Run("should return the buildings stored in the store", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)

			request, err := http.NewRequest("GET", "http://test/buildings", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

			actual := gateway.Buildings{}
			err = testutils.Read(res, &actual)
			if assert.NoError(t, err) {
				assert.Equal(t, buildings, actual)
			}
		})

		t.Run("should handle error returned by the store", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(nil, fmt.Errorf("unable to contact store"))

			request, err := http.NewRequest("GET", "http://test/buildings", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusInternalServerError, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "unable to contact store", msg)
			}
		})
	})

	t.Run("test POST /buildings", func(t *testing.T) {
		t.Run("should create building", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			building := gateway.Building{
				Lat:            1.2,
				Lan:            1.3,
				PhysicalEntity: gateway.PhysicalEntity{Name: "building-one", Description: "test-building"},
			}
			buildings := gateway.Buildings{}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().UpsertBuildings(gateway.Buildings{building.ID(): building}).Return(nil)

			data, _ := json.Marshal(building)

			request, err := http.NewRequest("POST", "http://test/buildings", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusCreated, res.StatusCode)
		})

		t.Run("should handle validation error if any", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			building := gateway.Building{
				Lat:            1.2,
				Lan:            1.3,
				PhysicalEntity: gateway.PhysicalEntity{Name: "one", Description: "test-building"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)

			data, _ := json.Marshal(building)

			request, err := http.NewRequest("POST", "http://test/buildings", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusBadRequest, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "name: the length must be between 5 and 50.", msg)
			}
		})

		t.Run("should handle error returned by store when fetching existing buildings", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			building := gateway.Building{
				Lat:            1.2,
				Lan:            1.3,
				PhysicalEntity: gateway.PhysicalEntity{Name: "building-one", Description: "test-building"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(nil, fmt.Errorf("unable to contact store"))

			data, _ := json.Marshal(building)

			request, err := http.NewRequest("POST", "http://test/buildings", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusInternalServerError, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "unable to contact store", msg)
			}
		})

		t.Run("should handle error returned by store when saving building", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			building := gateway.Building{
				Lat:            1.2,
				Lan:            1.3,
				PhysicalEntity: gateway.PhysicalEntity{Name: "building-one", Description: "test-building"},
			}
			buildings := gateway.Buildings{}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().UpsertBuildings(gateway.Buildings{building.ID(): building}).Return(fmt.Errorf("unable to save"))

			data, _ := json.Marshal(building)

			request, err := http.NewRequest("POST", "http://test/buildings", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusInternalServerError, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "unable to save", msg)
			}
		})
	})

	t.Run("test GET /buildings/:building-id", func(t *testing.T) {

		building := gateway.Building{
			Lat:            1.2,
			Lan:            1.3,
			PhysicalEntity: gateway.PhysicalEntity{Name: "building one", Description: "test building"},
		}

		buildings := gateway.Buildings{building.ID(): building}

		t.Run("should return the building", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)

			request, err := http.NewRequest("GET", "http://test/buildings/building-one", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

			actual := gateway.Building{}
			err = testutils.Read(res, &actual)
			if assert.NoError(t, err) {
				assert.Equal(t, building, actual)
			}
		})

		t.Run("should return 404 if building is not available", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)

			request, err := http.NewRequest("GET", "http://test/buildings/building-two", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, 404, res.StatusCode)
		})

		t.Run("should handle error returned by the store", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(nil, fmt.Errorf("unable to contact store"))

			request, err := http.NewRequest("GET", "http://test/buildings/one", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusInternalServerError, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "unable to contact store", msg)
			}
		})
	})

	t.Run("test PUT /buildings/:building-id", func(t *testing.T) {

		building := gateway.Building{
			Lat:            1.2,
			Lan:            1.3,
			PhysicalEntity: gateway.PhysicalEntity{Name: "building one", Description: "test building"},
		}

		buildings := gateway.Buildings{building.ID(): building}

		t.Run("should update building", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			building := gateway.Building{
				Lat:            1.2,
				Lan:            1.3,
				PhysicalEntity: gateway.PhysicalEntity{Name: "building-one", Description: "updated description"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().UpsertBuilding(building).Return(nil)

			data, _ := json.Marshal(building)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusOK, res.StatusCode)
		})

		t.Run("should return 404 if building is not available", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-two", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, 404, res.StatusCode)
		})

		t.Run("should change handle validation error if any", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			building := gateway.Building{
				Lat:            1.2,
				PhysicalEntity: gateway.PhysicalEntity{Name: "one", Description: "updated description"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)

			data, _ := json.Marshal(building)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusBadRequest, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "lan: cannot be blank; name: the length must be between 5 and 50.", msg)
			}
		})

		t.Run("should handle error returned by store when fetching existing buildings", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			building := gateway.Building{
				Lat:            1.2,
				Lan:            1.3,
				PhysicalEntity: gateway.PhysicalEntity{Name: "building-one", Description: "test-building"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(nil, fmt.Errorf("unable to contact store"))

			data, _ := json.Marshal(building)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusInternalServerError, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "unable to contact store", msg)
			}
		})

		t.Run("should handle error returned by store when saving building", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			building := gateway.Building{
				Lat:            1.2,
				Lan:            1.3,
				PhysicalEntity: gateway.PhysicalEntity{Name: "building-one", Description: "updated description"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().UpsertBuilding(building).Return(fmt.Errorf("unable to save"))

			data, _ := json.Marshal(building)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusInternalServerError, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "unable to save", msg)
			}
		})
	})

	t.Run("test DELETE /buildings/:building-id", func(t *testing.T) {

		building := gateway.Building{
			Lat:            1.2,
			Lan:            1.3,
			PhysicalEntity: gateway.PhysicalEntity{Name: "building one", Description: "test building"},
		}

		buildings := gateway.Buildings{building.ID(): building}

		t.Run("should delete building", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().DeleteBuilding(building).Return(nil)

			data, _ := json.Marshal(building)

			request, err := http.NewRequest("DELETE", "http://test/buildings/building-one", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusOK, res.StatusCode)
		})

		t.Run("should return 404 if building is not available", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-two", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, 404, res.StatusCode)
		})

		t.Run("should handle error returned by store when fetching existing buildings", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(nil, fmt.Errorf("unable to contact store"))

			data, _ := json.Marshal(building)

			request, err := http.NewRequest("DELETE", "http://test/buildings/building-one", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusInternalServerError, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "unable to contact store", msg)
			}
		})

		t.Run("should handle error returned by store when deleting building", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().DeleteBuilding(building).Return(fmt.Errorf("unable to delete"))

			data, _ := json.Marshal(building)

			request, err := http.NewRequest("DELETE", "http://test/buildings/building-one", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusInternalServerError, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "unable to delete", msg)
			}
		})
	})
}
