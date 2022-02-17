package store

import (
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"path"
)

const (
	buildingsBasePath = "buildings"
)

// Buildings returns all the Buildings from store
func (ps PersistentStore) Buildings() (gateway.Buildings, error) {
	value, err := ps.get(ps.buildingsRootPath(), gateway.Buildings{})
	if err != nil {
		return nil, err
	}
	return gateway.NewBuildings(value)
}

// UpsertBuildings creates or updates Buildings in store
func (ps PersistentStore) UpsertBuildings(buildings gateway.Buildings) error {
	return ps.putJSON(ps.buildingsRootPath(), buildings)
}

// UpsertBuilding creates or updates Building in store
func (ps PersistentStore) UpsertBuilding(building gateway.Building) error {
	buildings, err := ps.Buildings()
	if err != nil {
		return err
	}
	buildings[building.ID()] = building
	return ps.putJSON(ps.buildingsRootPath(), buildings)
}

// DeleteBuilding deletes the building and nested path from store
func (ps PersistentStore) DeleteBuilding(building gateway.Building) error {
	buildings, err := ps.Buildings()
	if err != nil {
		return err
	}

	delete(buildings, building.ID())
	err = ps.putJSON(ps.buildingsRootPath(), buildings)
	if err != nil {
		return err
	}

	return ps.safeDelete(ps.buildingRootPath(building))
}

func (ps PersistentStore) buildingsRootPath() string {
	return path.Join(ps.path, buildingsBasePath)
}

func (ps PersistentStore) buildingRootPath(building gateway.Entity) string {
	return path.Join(ps.path, building.ID())
}
