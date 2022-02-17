package store_test

import (
	"encoding/json"
	"fmt"
	libKVStore "github.com/abronan/valkeyrie/store"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	mockKVStore "gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/internal/mocks/valkeyrie/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/testutils"
	"testing"
)

func TestPersistentStore_Rooms(t *testing.T) {
	t.Run("should return rooms", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := testutils.NewFloor("floor-one")
		expectedRooms := gateway.Rooms{"room-one": gateway.Room{
			Direction: gateway.DirectionNorth,
			Floor:     building,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-one",
				Description: "test room",
			},
		}}
		data, _ := json.Marshal(expectedRooms)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floor-one/rooms", nil).Return(&libKVStore.KVPair{Value: data}, nil)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)

		actual, err := persistentStore.Rooms(building)

		assert.NoError(t, err)
		if !cmp.Equal(expectedRooms, actual) {
			assert.Fail(t, cmp.Diff(expectedRooms, actual))
		}
	})

	t.Run("should handle error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floor-one/rooms", nil).Return(nil, fmt.Errorf("store unavailable"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)

		actual, err := persistentStore.Rooms(testutils.NewFloor("floor-one"))

		if assert.Error(t, err) {
			assert.Equal(t, "store unavailable", err.Error())
		}
		assert.Nil(t, actual)
	})
}

func TestPersistentStore_UpsertRoom(t *testing.T) {
	t.Run("should add room to existing rooms", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := testutils.NewFloor("floor-one")
		room := gateway.Room{
			Direction: gateway.DirectionWest,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-one",
				Description: "test room",
			},
		}
		rooms := gateway.Rooms{"room-one": room}
		data, _ := json.Marshal(rooms)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floor-one/rooms", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/building-one/floor-one/rooms", gomock.Any(), nil).DoAndReturn(
			func(key string, data []byte, options *libKVStore.WriteOptions) error {
				actual := gateway.Rooms{}
				err := json.Unmarshal(data, &actual)
				if err != nil {
					return err
				}

				expected := gateway.Rooms{
					"room-one": room,
					"room-two": gateway.Room{
						Direction: gateway.DirectionWest,
						PhysicalEntity: gateway.PhysicalEntity{
							Name:        "room-two",
							Description: "test room",
						},
					},
				}
				if !cmp.Equal(expected, actual) {
					assert.Fail(t, cmp.Diff(expected, actual))
				}
				return nil
			},
		)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		newRoom := gateway.Room{
			Direction: gateway.DirectionWest,
			Floor:     building,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-two",
				Description: "test room",
			},
		}
		err := persistentStore.UpsertRoom(newRoom)
		assert.NoError(t, err)
	})

	t.Run("should handle error when saving rooms", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		room := gateway.Room{
			Direction: 1,
			Floor:     testutils.NewFloor("floor-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-one",
				Description: "test room",
			},
		}
		rooms := gateway.Rooms{}
		data, _ := json.Marshal(rooms)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floor-one/rooms", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/building-one/floor-one/rooms", gomock.Any(), nil).Return(fmt.Errorf("unable to save"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertRoom(room)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})

	t.Run("should handle error when fetching rooms", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		room := gateway.Room{
			Direction: 1,
			Floor:     testutils.NewFloor("floor-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-one",
				Description: "test room",
			},
		}

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floor-one/rooms", nil).Return(nil, fmt.Errorf("unable to fetch rooms"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertRoom(room)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to fetch rooms", err.Error())
		}
	})
}

func TestPersistentStore_UpsertRooms(t *testing.T) {
	t.Run("should add room to existing rooms", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		room := gateway.Room{
			Direction: gateway.DirectionSouth,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-one",
				Description: "test room",
			},
		}
		rooms := gateway.Rooms{"existing": room}

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Put("dwarka/building-one/floor-one/rooms", gomock.Any(), nil).DoAndReturn(
			func(key string, data []byte, options *libKVStore.WriteOptions) error {
				actual := gateway.Rooms{}
				err := json.Unmarshal(data, &actual)
				if err != nil {
					return err
				}

				assert.Equal(t, gateway.Rooms{
					"existing": room,
				}, actual)
				return nil
			},
		)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertRooms(testutils.NewFloor("floor-one"), rooms)
		assert.NoError(t, err)
	})

	t.Run("should handle error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		room := gateway.Room{
			Direction: 1,
			Floor:     testutils.NewFloor("floor-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-one",
				Description: "test room",
			},
		}
		rooms := gateway.Rooms{"existing": room}

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Put("dwarka/building-one/floor-one/rooms", gomock.Any(), nil).Return(fmt.Errorf("unable to save"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertRooms(testutils.NewFloor("floor-one"), rooms)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})
}

func TestPersistentStore_DeleteRoom(t *testing.T) {

	t.Run("should delete room to from rooms", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		room := gateway.Room{
			Direction: 1,
			Floor:     testutils.NewFloor("floor-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "existing",
				Description: "test room",
			},
		}
		rooms := gateway.Rooms{"existing": room}
		data, _ := json.Marshal(rooms)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floor-one/rooms", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().DeleteTree("dwarka/building-one/floor-one/existing").Return(nil)
		mockStore.EXPECT().Put("dwarka/building-one/floor-one/rooms", gomock.Any(), nil).DoAndReturn(
			func(key string, data []byte, options *libKVStore.WriteOptions) error {
				actual := gateway.Rooms{}
				err := json.Unmarshal(data, &actual)
				if err != nil {
					return err
				}

				assert.Equal(t, gateway.Rooms{}, actual)
				return nil
			},
		)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteRoom(room)
		assert.NoError(t, err)
	})

	t.Run("should handle error when deleting rooms", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		room := gateway.Room{
			Direction: 1,
			Floor:     testutils.NewFloor("floor-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-one",
				Description: "test room",
			},
		}
		rooms := gateway.Rooms{}
		data, _ := json.Marshal(rooms)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floor-one/rooms", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/building-one/floor-one/rooms", gomock.Any(), nil).Return(nil)
		mockStore.EXPECT().DeleteTree("dwarka/building-one/floor-one/room-one").Return(fmt.Errorf("unable to delete"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteRoom(room)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to delete", err.Error())
		}
	})

	t.Run("should handle error when saving rooms", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		room := gateway.Room{
			Direction: 1,
			Floor:     testutils.NewFloor("floor-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-one",
				Description: "test room",
			},
		}
		rooms := gateway.Rooms{}
		data, _ := json.Marshal(rooms)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floor-one/rooms", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/building-one/floor-one/rooms", gomock.Any(), nil).Return(fmt.Errorf("unable to save"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteRoom(room)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})

	t.Run("should handle error when fetching existing rooms", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		room := gateway.Room{
			Direction: 1,
			Floor:     testutils.NewFloor("floor-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-one",
				Description: "test room",
			},
		}

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floor-one/rooms", nil).Return(nil, fmt.Errorf("unable to fetch rooms"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteRoom(room)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to fetch rooms", err.Error())
		}
	})
}
