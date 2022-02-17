package store

import (
	"encoding/json"
	"fmt"

	"github.com/kvtools/valkeyrie"
	"github.com/kvtools/valkeyrie/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
)

//go:generate $PWD/scripts/mockgen $PWD/pkg/store/store.go $PWD/pkg/internal/mocks/store/store.go mockStore
//go:generate $PWD/scripts/mockgen $PWD/vendor/github.com/kvtools/valkeyrie/store/store.go $PWD/pkg/internal/mocks/valkeyrie/store/store.go mockKVStore

// Store represents the backend for dwarka
// Each store should support every call listed
type Store interface {
	Buildings() (gateway.Buildings, error)
	UpsertBuildings(buildings gateway.Buildings) error
	UpsertBuilding(building gateway.Building) error
	DeleteBuilding(building gateway.Building) error
	Floors(building gateway.Entity) (gateway.Floors, error)
	UpsertFloors(building gateway.Entity, floors gateway.Floors) error
	UpsertFloor(floor gateway.Floor) error
	DeleteFloor(floor gateway.Floor) error
	Rooms(floor gateway.Floor) (gateway.Rooms, error)
	UpsertRooms(floor gateway.Floor, rooms gateway.Rooms) error
	UpsertRoom(room gateway.Room) error
	DeleteRoom(room gateway.Room) error
	Uptime() (gateway.Status, error)
	RefreshUptime() error
}

// NotFound is thrown when the key is not found in the store during a Get operation
type NotFound string

// Error returns the underlying error as string
func (err NotFound) Error() string {
	return string(err)
}

// PersistentStore is a persistent implementation for Store
// the data is persisted in one of the kv store supported by libkv
type PersistentStore struct {
	path    string
	kvStore store.Store
}

func (ps PersistentStore) get(path string, defaultValue interface{}) ([]byte, error) {
	kv, err := ps.kvStore.Get(path, nil)
	if err == store.ErrKeyNotFound {
		return json.Marshal(defaultValue)
	} else if err != nil {
		return nil, err
	}
	return kv.Value, nil
}

func (ps PersistentStore) put(path string, value []byte) error {
	return ps.kvStore.Put(path, value, nil)
}

func (ps PersistentStore) putJSON(path string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ps.put(path, data)
}

func (ps PersistentStore) safeDelete(path string) error {
	err := ps.kvStore.DeleteTree(path)
	if err != nil && err != store.ErrKeyNotFound {
		return err
	}
	return nil
}

// NewPersistentStore returns a instance of PersistentStore
func NewPersistentStore(path string, store store.Store) Store {
	return &PersistentStore{path: path, kvStore: store}
}

// NewStore return libkv/store.PersistentStore with necessary defaults
func NewStore(basePath string, backend string, bucketName string, addrs ...string) (Store, error) {
	s, err := valkeyrie.NewStore(store.Backend(backend), addrs, &store.Config{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create backend store, reason: %v", err)
	}
	return NewPersistentStore(basePath, s), nil
}
