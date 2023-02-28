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

func TestFloors(t *testing.T) {

	t.Run("test GET /buildings/:building-id/floors", func(t *testing.T) {

		floor := gateway.Floor{
			Level:          1,
			PhysicalEntity: gateway.PhysicalEntity{Name: "floor one", Description: "test floor"},
		}
		floors := gateway.Floors{floor.ID(): floor}

		t.Run("should return the floors stored in the store", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			buildings, building := testutils.NewBuildings("building-one")

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)

			request, err := http.NewRequest("GET", "http://test/buildings/building-one/floors", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

			actual := gateway.Floors{}
			err = testutils.Read(res, &actual)
			if assert.NoError(t, err) {
				assert.Equal(t, floors, actual)
			}
		})

		t.Run("should handle error returned by the store", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			buildings, building := testutils.NewBuildings("building-one")

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(nil, fmt.Errorf("unable to contact store"))

			request, err := http.NewRequest("GET", "http://test/buildings/building-one/floors", nil)
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

	t.Run("test POST /buildings/:building-id/floors", func(t *testing.T) {

		t.Run("should create floor", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			floor := gateway.Floor{
				Level:          1,
				PhysicalEntity: gateway.PhysicalEntity{Name: "floor-one", Description: "test-floor"},
			}
			floors := gateway.Floors{}
			buildings, building := testutils.NewBuildings("building-one")

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().UpsertFloors(building, gomock.Any()).DoAndReturn(
				func(_ gateway.Building, actualFloors gateway.Floors) error {
					expectedFloors := gateway.Floors{floor.ID(): testutils.AssociateFloorToBuilding(building, floor)}
					assert.Equal(t, expectedFloors, actualFloors)
					return nil
				},
			)

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("POST", "http://test/buildings/building-one/floors", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusCreated, res.StatusCode)
		})

		t.Run("should return 409 if floor already exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			floor := gateway.Floor{
				Level:          1,
				PhysicalEntity: gateway.PhysicalEntity{Name: "floor-one", Description: "test-floor"},
			}
			floors := gateway.Floors{floor.ID(): floor}
			buildings, building := testutils.NewBuildings("building-one")

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("POST", "http://test/buildings/building-one/floors", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusConflict, res.StatusCode)
		})

		t.Run("should handle validation error if any", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			floor := gateway.Floor{
				Level:          1,
				PhysicalEntity: gateway.PhysicalEntity{Name: "", Description: "test-floor"},
			}
			buildings, _ := testutils.NewBuildings("building-one")

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("POST", "http://test/buildings/building-one/floors", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusBadRequest, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "name: cannot be blank.", msg)
			}
		})

		t.Run("should handle error returned by store when fetching existing floors", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			floor := gateway.Floor{
				Level:          1,
				PhysicalEntity: gateway.PhysicalEntity{Name: "floor-one", Description: "test-floor"},
			}
			buildings, _ := testutils.NewBuildings("building-one")

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(testutils.NewBuilding("building-one")).Return(nil, fmt.Errorf("unable to contact store"))

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("POST", "http://test/buildings/building-one/floors", bytes.NewReader(data))
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

		t.Run("should handle error returned by store when saving floor", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			floor := gateway.Floor{
				Level:          1,
				PhysicalEntity: gateway.PhysicalEntity{Name: "floor-one", Description: "test-floor"},
			}
			buildings, building := testutils.NewBuildings("building-one")
			floors := gateway.Floors{}

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().UpsertFloors(building, gomock.Any()).Return(fmt.Errorf("unable to save"))

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("POST", "http://test/buildings/building-one/floors", bytes.NewReader(data))
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

	t.Run("test GET /floors/:floor-id", func(t *testing.T) {
		buildings, building := testutils.NewBuildings("building-one")
		floor := gateway.Floor{
			Level:          1,
			PhysicalEntity: gateway.PhysicalEntity{Name: "floor one", Description: "test floor"},
		}

		floors := gateway.Floors{floor.ID(): floor}

		t.Run("should return the floor", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)

			request, err := http.NewRequest("GET", "http://test/buildings/building-one/floors/floor-one", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

			actual := gateway.Floor{}
			err = testutils.Read(res, &actual)
			if assert.NoError(t, err) {
				assert.Equal(t, floor, actual)
			}
		})

		t.Run("should return 404 if floor is not available", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			building := testutils.NewBuilding("building-one")
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)

			request, err := http.NewRequest("GET", "http://test/buildings/building-one/floors/floor-two", nil)
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

			building := testutils.NewBuilding("building-one")
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(nil, fmt.Errorf("unable to contact store"))

			request, err := http.NewRequest("GET", "http://test/buildings/building-one/floors/one", nil)
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

	t.Run("test PUT /floors/:floor-id", func(t *testing.T) {
		buildings, building := testutils.NewBuildings("building-one")
		floor := gateway.Floor{
			Level:          1,
			PhysicalEntity: gateway.PhysicalEntity{Name: "floor one", Description: "test floor"},
		}

		floors := gateway.Floors{floor.ID(): floor}

		t.Run("should update floor", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			floor := gateway.Floor{
				Level:          1,
				Building:       building,
				PhysicalEntity: gateway.PhysicalEntity{Name: "floor-one", Description: "updated description"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().UpsertFloor(floor).Return(nil)

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one/floors/floor-one", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusOK, res.StatusCode)
		})

		t.Run("should return 404 if floor is not available", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one/floors/floor-two", nil)
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

			floor := gateway.Floor{
				PhysicalEntity: gateway.PhysicalEntity{Name: "one", Description: "updated description"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one/floors/floor-one", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusBadRequest, res.StatusCode)

			msg, err := testutils.ReadError(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "level: cannot be blank; name: the length must be between 5 and 50.", msg)
			}
		})

		t.Run("should handle error returned by store when fetching existing floors", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			floor := gateway.Floor{
				Level:          1,
				PhysicalEntity: gateway.PhysicalEntity{Name: "floor-one", Description: "test-floor"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(nil, fmt.Errorf("unable to contact store"))

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one/floors/floor-one", bytes.NewReader(data))
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

		t.Run("should handle error returned by store when saving floor", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			floor := gateway.Floor{
				Level:          1,
				Building:       building,
				PhysicalEntity: gateway.PhysicalEntity{Name: "floor-one", Description: "updated description"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().UpsertFloor(floor).Return(fmt.Errorf("unable to save"))

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one/floors/floor-one", bytes.NewReader(data))
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

	t.Run("test DELETE /floors/:floor-id", func(t *testing.T) {
		buildings, building := testutils.NewBuildings("building-one")
		floor := gateway.Floor{
			Level:          1,
			PhysicalEntity: gateway.PhysicalEntity{Name: "floor one", Description: "test floor"},
		}

		floors := gateway.Floors{floor.ID(): floor}

		t.Run("should delete floor", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().DeleteFloor(floor).Return(nil)

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("DELETE", "http://test/buildings/building-one/floors/floor-one", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusOK, res.StatusCode)
		})

		t.Run("should return 404 if floor is not available", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one/floors/floor-two", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, 404, res.StatusCode)
		})

		t.Run("should handle error returned by store when fetching existing floors", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			building := testutils.NewBuilding("building-one")
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(nil, fmt.Errorf("unable to contact store"))

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("DELETE", "http://test/buildings/building-one/floors/floor-one", bytes.NewReader(data))
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

		t.Run("should handle error returned by store when deleting floor", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().DeleteFloor(floor).Return(fmt.Errorf("unable to delete"))

			data, _ := json.Marshal(floor)

			request, err := http.NewRequest("DELETE", "http://test/buildings/building-one/floors/floor-one", bytes.NewReader(data))
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
