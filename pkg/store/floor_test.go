package store_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	libKVStore "github.com/kvtools/valkeyrie/store"
	"github.com/stretchr/testify/assert"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	mockKVStore "gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/internal/mocks/valkeyrie/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/testutils"
)

func TestPersistentStore_Floors(t *testing.T) {
	t.Run("should return floors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := testutils.NewBuilding("building-one")
		expectedFloors := gateway.Floors{"floor-one": gateway.Floor{
			Level:    1,
			Building: building,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-one",
				Description: "test floor",
			},
		}}
		data, _ := json.Marshal(expectedFloors)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floors", nil).Return(&libKVStore.KVPair{Value: data}, nil)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)

		actual, err := persistentStore.Floors(building)

		assert.NoError(t, err)
		if !cmp.Equal(expectedFloors, actual) {
			assert.Fail(t, cmp.Diff(expectedFloors, actual))
		}
	})

	t.Run("should handle error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floors", nil).Return(nil, fmt.Errorf("store unavailable"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)

		actual, err := persistentStore.Floors(testutils.NewBuilding("building-one"))

		if assert.Error(t, err) {
			assert.Equal(t, "store unavailable", err.Error())
		}
		assert.Nil(t, actual)
	})
}

func TestPersistentStore_UpsertFloor(t *testing.T) {
	t.Run("should add floor to existing floors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := testutils.NewBuilding("building-one")
		floor := gateway.Floor{
			Level: 1,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-one",
				Description: "test floor",
			},
		}
		floors := gateway.Floors{"floor-one": floor}
		data, _ := json.Marshal(floors)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floors", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/building-one/floors", gomock.Any(), nil).DoAndReturn(
			func(key string, data []byte, options *libKVStore.WriteOptions) error {
				actual := gateway.Floors{}
				err := json.Unmarshal(data, &actual)
				if err != nil {
					return err
				}

				expected := gateway.Floors{
					"floor-one": floor,
					"floor-two": gateway.Floor{
						Level: 1,
						PhysicalEntity: gateway.PhysicalEntity{
							Name:        "floor-two",
							Description: "test floor",
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
		newFloor := gateway.Floor{
			Level:    1,
			Building: building,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-two",
				Description: "test floor",
			},
		}
		err := persistentStore.UpsertFloor(newFloor)
		assert.NoError(t, err)
	})

	t.Run("should handle error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		floor := gateway.Floor{
			Level:    1,
			Building: testutils.NewBuilding("building-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-one",
				Description: "test floor",
			},
		}
		floors := gateway.Floors{}
		data, _ := json.Marshal(floors)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floors", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/building-one/floors", gomock.Any(), nil).Return(fmt.Errorf("unable to save"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertFloor(floor)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})

	t.Run("should handle error when getting existing floors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		floor := gateway.Floor{
			Level:    1,
			Building: testutils.NewBuilding("building-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-one",
				Description: "test floor",
			},
		}
		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floors", nil).Return(nil, fmt.Errorf("unable to get floors"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertFloor(floor)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to get floors", err.Error())
		}
	})
}

func TestPersistentStore_UpsertFloors(t *testing.T) {
	t.Run("should add floor to existing floors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		floor := gateway.Floor{
			Level: 1,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-one",
				Description: "test floor",
			},
		}
		floors := gateway.Floors{"existing": floor}

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Put("dwarka/building-one/floors", gomock.Any(), nil).DoAndReturn(
			func(key string, data []byte, options *libKVStore.WriteOptions) error {
				actual := gateway.Floors{}
				err := json.Unmarshal(data, &actual)
				if err != nil {
					return err
				}

				assert.Equal(t, gateway.Floors{
					"existing": floor,
				}, actual)
				return nil
			},
		)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertFloors(testutils.NewBuilding("building-one"), floors)
		assert.NoError(t, err)
	})

	t.Run("should handle error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		floor := gateway.Floor{
			Level:    1,
			Building: testutils.NewBuilding("building-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-one",
				Description: "test floor",
			},
		}
		floors := gateway.Floors{"existing": floor}

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Put("dwarka/building-one/floors", gomock.Any(), nil).Return(fmt.Errorf("unable to save"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertFloors(testutils.NewBuilding("building-one"), floors)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})
}

func TestPersistentStore_DeleteFloor(t *testing.T) {

	t.Run("should delete floor to from floors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		floor := gateway.Floor{
			Level:    1,
			Building: testutils.NewBuilding("building-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "existing",
				Description: "test floor",
			},
		}
		floors := gateway.Floors{"existing": floor}
		data, _ := json.Marshal(floors)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floors", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().DeleteTree("dwarka/building-one/existing").Return(nil)
		mockStore.EXPECT().Put("dwarka/building-one/floors", gomock.Any(), nil).DoAndReturn(
			func(key string, data []byte, options *libKVStore.WriteOptions) error {
				actual := gateway.Floors{}
				err := json.Unmarshal(data, &actual)
				if err != nil {
					return err
				}

				assert.Equal(t, gateway.Floors{}, actual)
				return nil
			},
		)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteFloor(floor)
		assert.NoError(t, err)
	})

	t.Run("should handle error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		floor := gateway.Floor{
			Level:    1,
			Building: testutils.NewBuilding("building-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-one",
				Description: "test floor",
			},
		}
		floors := gateway.Floors{}
		data, _ := json.Marshal(floors)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floors", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/building-one/floors", gomock.Any(), nil).Return(nil)
		mockStore.EXPECT().DeleteTree("dwarka/building-one/floor-one").Return(fmt.Errorf("unable to delete"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteFloor(floor)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to delete", err.Error())
		}
	})

	t.Run("should handle error when fetching existing floors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		floor := gateway.Floor{
			Level:    1,
			Building: testutils.NewBuilding("building-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-one",
				Description: "test floor",
			},
		}

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floors", nil).Return(nil, fmt.Errorf("unable to get floors"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteFloor(floor)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to get floors", err.Error())
		}
	})

	t.Run("should handle error when updating existing floors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		floor := gateway.Floor{
			Level:    1,
			Building: testutils.NewBuilding("building-one"),
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-one",
				Description: "test floor",
			},
		}
		floors := gateway.Floors{}
		data, _ := json.Marshal(floors)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/building-one/floors", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/building-one/floors", gomock.Any(), nil).Return(fmt.Errorf("unable to update floors"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteFloor(floor)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to update floors", err.Error())
		}
	})
}
