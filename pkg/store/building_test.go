package store_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	libKVStore "github.com/kvtools/valkeyrie/store"
	"github.com/stretchr/testify/assert"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	mockKVStore "gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/internal/mocks/valkeyrie/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
)

func TestPersistentStore_Buildings(t *testing.T) {
	t.Run("should return dwarka/buildings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		expectedBuilds := gateway.Buildings{"one": gateway.Building{
			Lat: 1.2,
			Lan: 1.3,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "building-one",
				Description: "test building",
			},
		}}
		data, _ := json.Marshal(expectedBuilds)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/buildings", nil).Return(&libKVStore.KVPair{Value: data}, nil)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)

		actual, err := persistentStore.Buildings()

		assert.NoError(t, err)
		assert.Equal(t, expectedBuilds, actual)
	})

	t.Run("should handle error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/buildings", nil).Return(nil, fmt.Errorf("store unavailable"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)

		actual, err := persistentStore.Buildings()

		if assert.Error(t, err) {
			assert.Equal(t, "store unavailable", err.Error())
		}
		assert.Nil(t, actual)
	})
}

func TestPersistentStore_UpsertBuilding(t *testing.T) {
	t.Run("should add building to existing dwarka/buildings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := gateway.Building{
			Lat: 1.2,
			Lan: 1.3,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "building-one",
				Description: "test building",
			},
		}
		buildings := gateway.Buildings{"existing": building}
		data, _ := json.Marshal(buildings)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/buildings", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/buildings", gomock.Any(), nil).DoAndReturn(
			func(key string, data []byte, options *libKVStore.WriteOptions) error {
				actual := gateway.Buildings{}
				err := json.Unmarshal(data, &actual)
				if err != nil {
					return err
				}

				assert.Equal(t, gateway.Buildings{
					building.ID(): building,
					"existing":    building,
				}, actual)
				return nil
			},
		)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertBuilding(building)
		assert.NoError(t, err)
	})

	t.Run("should handle error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := gateway.Building{
			Lat: 1.2,
			Lan: 1.3,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "building-one",
				Description: "test building",
			},
		}
		buildings := gateway.Buildings{}
		data, _ := json.Marshal(buildings)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/buildings", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/buildings", gomock.Any(), nil).Return(fmt.Errorf("unable to save"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertBuilding(building)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})

	t.Run("should handle error when fetching buildings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := gateway.Building{
			Lat: 1.2,
			Lan: 1.3,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "building-one",
				Description: "test building",
			},
		}

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/buildings", nil).Return(nil, fmt.Errorf("unable to get buildings"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertBuilding(building)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to get buildings", err.Error())
		}
	})
}

func TestPersistentStore_UpsertBuildings(t *testing.T) {
	t.Run("should add building to existing dwarka/buildings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := gateway.Building{
			Lat: 1.2,
			Lan: 1.3,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "building-one",
				Description: "test building",
			},
		}
		buildings := gateway.Buildings{"existing": building}

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Put("dwarka/buildings", gomock.Any(), nil).DoAndReturn(
			func(key string, data []byte, options *libKVStore.WriteOptions) error {
				actual := gateway.Buildings{}
				err := json.Unmarshal(data, &actual)
				if err != nil {
					return err
				}

				assert.Equal(t, gateway.Buildings{
					"existing": building,
				}, actual)
				return nil
			},
		)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertBuildings(buildings)
		assert.NoError(t, err)
	})

	t.Run("should handle error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := gateway.Building{
			Lat: 1.2,
			Lan: 1.3,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "building-one",
				Description: "test building",
			},
		}
		buildings := gateway.Buildings{"existing": building}

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Put("dwarka/buildings", gomock.Any(), nil).Return(fmt.Errorf("unable to save"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.UpsertBuildings(buildings)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})
}

func TestPersistentStore_DeleteBuilding(t *testing.T) {

	t.Run("should delete building to from dwarka/buildings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := gateway.Building{
			Lat: 1.2,
			Lan: 1.3,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "existing",
				Description: "test building",
			},
		}
		buildings := gateway.Buildings{"existing": building}
		data, _ := json.Marshal(buildings)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/buildings", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().DeleteTree("dwarka/existing").Return(nil)
		mockStore.EXPECT().Put("dwarka/buildings", gomock.Any(), nil).DoAndReturn(
			func(key string, data []byte, options *libKVStore.WriteOptions) error {
				actual := gateway.Buildings{}
				err := json.Unmarshal(data, &actual)
				if err != nil {
					return err
				}

				assert.Equal(t, gateway.Buildings{}, actual)
				return nil
			},
		)

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteBuilding(building)
		assert.NoError(t, err)
	})

	t.Run("should handle error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := gateway.Building{
			Lat: 1.2,
			Lan: 1.3,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "building-one",
				Description: "test building",
			},
		}
		buildings := gateway.Buildings{}
		data, _ := json.Marshal(buildings)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/buildings", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/buildings", gomock.Any(), nil).Return(nil)
		mockStore.EXPECT().DeleteTree("dwarka/building-one").Return(fmt.Errorf("unable to delete"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteBuilding(building)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to delete", err.Error())
		}
	})

	t.Run("should handle error when get existing buildings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := gateway.Building{
			Lat: 1.2,
			Lan: 1.3,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "building-one",
				Description: "test building",
			},
		}

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/buildings", nil).Return(nil, fmt.Errorf("unable to get buildings"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteBuilding(building)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to get buildings", err.Error())
		}
	})

	t.Run("should handle error when updating buildings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		building := gateway.Building{
			Lat: 1.2,
			Lan: 1.3,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "building-one",
				Description: "test building",
			},
		}
		buildings := gateway.Buildings{}
		data, _ := json.Marshal(buildings)

		mockStore := mockKVStore.NewMockStore(ctrl)
		mockStore.EXPECT().Get("dwarka/buildings", nil).Return(&libKVStore.KVPair{Value: data}, nil)
		mockStore.EXPECT().Put("dwarka/buildings", gomock.Any(), nil).Return(fmt.Errorf("unable to update buildings"))

		persistentStore := store.NewPersistentStore("dwarka", mockStore)
		err := persistentStore.DeleteBuilding(building)
		if assert.Error(t, err) {
			assert.Equal(t, "unable to update buildings", err.Error())
		}
	})
}
