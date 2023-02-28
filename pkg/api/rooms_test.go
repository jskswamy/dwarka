package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/api/view"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	mockStore "gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/internal/mocks/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/testutils"
	"net/http"
	"testing"
)

func TestRooms(t *testing.T) {

	t.Run("test GET /buildings/:building-id/floors/floor-one/rooms", func(t *testing.T) {

		floors, floor := testutils.NewFloors("floor-one")
		room := gateway.Room{
			Direction:      gateway.DirectionNorth,
			PhysicalEntity: gateway.PhysicalEntity{Name: "room one", Description: "test floor"},
		}
		rooms := gateway.Rooms{room.ID(): room}
		expected := view.NewRooms(rooms)

		t.Run("should return the rooms stored in the store", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			buildings, building := testutils.NewBuildings("building-one")

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)

			request, err := http.NewRequest("GET", "http://test/buildings/building-one/floors/floor-one/rooms", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

			var actual []view.Room
			err = testutils.Read(res, &actual)
			if assert.NoError(t, err) {
				assert.Equal(t, expected, actual)
			}
		})

		t.Run("should handle error returned by the store", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			buildings, building := testutils.NewBuildings("building-one")

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(nil, fmt.Errorf("unable to contact store"))

			request, err := http.NewRequest("GET", "http://test/buildings/building-one/floors/floor-one/rooms", nil)
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

	t.Run("test POST /buildings/:building-id/floors/floor-one/rooms", func(t *testing.T) {

		floors, floor := testutils.NewFloors("floor-one")
		buildings, building := testutils.NewBuildings("building-one")
		room := gateway.Room{
			Direction:      gateway.DirectionNorth,
			PhysicalEntity: gateway.PhysicalEntity{Name: "room one", Description: "test floor"},
		}
		newRoom := view.Room{
			Direction:   "north",
			Name:        "room one",
			Description: "test floor",
		}

		t.Run("should create room", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rooms := gateway.Rooms{}

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)
			mockKVStore.EXPECT().UpsertRooms(floor, gomock.Any()).DoAndReturn(
				func(_ gateway.Floor, actualRooms gateway.Rooms) error {
					expectedRooms := gateway.Rooms{room.ID(): gateway.Room{
						Floor:          floor,
						Direction:      gateway.DirectionNorth,
						PhysicalEntity: room.PhysicalEntity,
					}}
					if !cmp.Equal(expectedRooms, actualRooms) {
						assert.Fail(t, cmp.Diff(expectedRooms, actualRooms))
					}
					return nil
				},
			)

			data, _ := json.Marshal(newRoom)

			request, err := http.NewRequest("POST", "http://test/buildings/building-one/floors/floor-one/rooms", bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusCreated, res.StatusCode)
		})

		t.Run("should return 409 if the room already exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rooms := gateway.Rooms{room.ID(): room}

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)

			data, _ := json.Marshal(newRoom)

			request, err := http.NewRequest("POST", "http://test/buildings/building-one/floors/floor-one/rooms", bytes.NewReader(data))
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

			invalidRoom := view.Room{
				Direction:   "north",
				Name:        "",
				Description: "test room",
			}

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)

			data, _ := json.Marshal(invalidRoom)

			request, err := http.NewRequest("POST", "http://test/buildings/building-one/floors/floor-one/rooms", bytes.NewReader(data))
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

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(testutils.NewBuilding("building-one")).Return(nil, fmt.Errorf("unable to contact store"))

			data, _ := json.Marshal(room)

			request, err := http.NewRequest("POST", "http://test/buildings/building-one/floors/floor-one/rooms", bytes.NewReader(data))
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

		t.Run("should handle error returned by store when fetching existing rooms", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(testutils.NewFloor("floor-one")).Return(nil, fmt.Errorf("unable to contact store"))

			data, _ := json.Marshal(newRoom)

			request, err := http.NewRequest("POST", "http://test/buildings/building-one/floors/floor-one/rooms", bytes.NewReader(data))
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

			rooms := gateway.Rooms{}

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)
			mockKVStore.EXPECT().UpsertRooms(floor, gomock.Any()).Return(fmt.Errorf("unable to save"))

			data, _ := json.Marshal(newRoom)

			request, err := http.NewRequest("POST", "http://test/buildings/building-one/floors/floor-one/rooms", bytes.NewReader(data))
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

	t.Run("test GET /buildings/:building-one/floors/:floor-id/rooms/:room-id", func(t *testing.T) {

		floors, floor := testutils.NewFloors("floor-one")
		buildings, building := testutils.NewBuildings("building-one")
		room := gateway.Room{
			Direction:      gateway.DirectionNorth,
			PhysicalEntity: gateway.PhysicalEntity{Name: "room one", Description: "test floor"},
		}
		rooms := gateway.Rooms{room.ID(): room}
		expected := view.NewRoom(room)

		t.Run("should return the floor", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)

			request, err := http.NewRequest("GET", "http://test/buildings/building-one/floors/floor-one/rooms/room-one", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

			actual := view.Room{}
			err = testutils.Read(res, &actual)
			if assert.NoError(t, err) {
				if !cmp.Equal(expected, actual) {
					assert.Fail(t, cmp.Diff(expected, actual))
				}
			}
		})

		t.Run("should return 404 if floor is not available", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			building := testutils.NewBuilding("building-one")
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)

			request, err := http.NewRequest("GET", "http://test/buildings/building-one/floors/floor-one/rooms/room-two", nil)
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
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(nil, fmt.Errorf("unable to contact store"))

			request, err := http.NewRequest("GET", "http://test/buildings/building-one/floors/floor-one/rooms/room-one", nil)
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

	t.Run("test PUT /buildings/:building-one/floors/:floor-id/rooms/:room-id", func(t *testing.T) {

		floors, floor := testutils.NewFloors("floor-one")
		buildings, building := testutils.NewBuildings("building-one")
		room := gateway.Room{
			Direction:      gateway.DirectionNorth,
			PhysicalEntity: gateway.PhysicalEntity{Name: "room one", Description: "test floor"},
		}
		rooms := gateway.Rooms{room.ID(): room}

		t.Run("should update room", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			updatedRoom := gateway.Room{
				Direction:      gateway.DirectionNorth,
				Floor:          floor,
				PhysicalEntity: gateway.PhysicalEntity{Name: "floor-one", Description: "updated description"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)
			mockKVStore.EXPECT().UpsertRoom(updatedRoom).Return(nil)

			data, _ := json.Marshal(view.NewRoom(updatedRoom))

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one/floors/floor-one/rooms/room-one", bytes.NewReader(data))
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
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one/floors/floor-one/rooms/room-two", nil)
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

			updatedRoom := gateway.Room{
				PhysicalEntity: gateway.PhysicalEntity{Name: "one", Description: "updated description"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)

			data, _ := json.Marshal(updatedRoom)

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

		t.Run("should handle error returned by store when fetching existing rooms", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			room := gateway.Room{
				Direction:      gateway.DirectionNorth,
				PhysicalEntity: gateway.PhysicalEntity{Name: "floor-one", Description: "test room"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(nil, fmt.Errorf("unable to contact store"))

			data, _ := json.Marshal(room)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one/floors/floor-one/rooms/room-one", bytes.NewReader(data))
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

			updatedRoom := gateway.Room{
				Direction:      gateway.DirectionNorth,
				Floor:          floor,
				PhysicalEntity: gateway.PhysicalEntity{Name: "floor-one", Description: "updated description"},
			}
			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)
			mockKVStore.EXPECT().UpsertRoom(updatedRoom).Return(fmt.Errorf("unable to save"))

			data, _ := json.Marshal(view.NewRoom(updatedRoom))

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one/floors/floor-one/rooms/room-one", bytes.NewReader(data))
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

	t.Run("test DELETE /buildings/:building-one/floors/:floor-id/rooms/:room-id", func(t *testing.T) {

		floors, floor := testutils.NewFloors("floor-one")
		buildings, building := testutils.NewBuildings("building-one")
		room := gateway.Room{
			Direction:      gateway.DirectionNorth,
			PhysicalEntity: gateway.PhysicalEntity{Name: "room one", Description: "test floor"},
		}
		rooms := gateway.Rooms{room.ID(): room}

		t.Run("should delete floor", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)
			mockKVStore.EXPECT().DeleteRoom(room).Return(nil)

			data, _ := json.Marshal(room)

			request, err := http.NewRequest("DELETE", "http://test/buildings/building-one/floors/floor-one/rooms/room-one", bytes.NewReader(data))
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
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)

			request, err := http.NewRequest("PUT", "http://test/buildings/building-one/floors/floor-one/rooms/room-two", nil)
			if err != nil {
				t.Error(err)
			}

			res, err := testutils.ServeHTTPRequest(mockKVStore, request)
			assert.NoError(t, err)
			assert.Equal(t, 404, res.StatusCode)
		})

		t.Run("should handle error returned by store when fetching existing rooms", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockKVStore := mockStore.NewMockStore(ctrl)
			mockKVStore.EXPECT().Buildings().Return(buildings, nil)
			mockKVStore.EXPECT().Floors(building).Return(floors, nil)
			mockKVStore.EXPECT().Rooms(floor).Return(nil, fmt.Errorf("unable to contact store"))

			data, _ := json.Marshal(room)

			request, err := http.NewRequest("DELETE", "http://test/buildings/building-one/floors/floor-one/rooms/room-one", bytes.NewReader(data))
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
			mockKVStore.EXPECT().Rooms(floor).Return(rooms, nil)
			mockKVStore.EXPECT().DeleteRoom(room).Return(fmt.Errorf("unable to delete"))

			data, _ := json.Marshal(room)

			request, err := http.NewRequest("DELETE", "http://test/buildings/building-one/floors/floor-one/rooms/room-one", bytes.NewReader(data))
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
