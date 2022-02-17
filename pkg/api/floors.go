package api

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/api/server"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
	"net/http"
	"path"
)

const (
	floorID       = "floor-id"
	floorUserKey  = "floor"
	floorsUserKey = "floors"
)

func floorPath() string {
	return path.Join(floorsBasePath(), fmt.Sprintf(":%s", floorID))
}

func floorsBasePath() string {
	return path.Join(buildingPath(), "floors")
}

func init() {
	floorFilters := &server.Filters{Before: []server.ResponseHandler{findAndLoadFloor}}
	buildingFilters := &server.Filters{Before: []server.ResponseHandler{findAndLoadBuilding}}
	AddRoute(
		server.NewRouteWithFilters("GET", floorsBasePath(), listFloorsHandler, buildingFilters),
		server.NewRouteWithFilters("POST", floorsBasePath(), createFloorHandler, buildingFilters),
		server.NewRouteWithFilters("GET", floorPath(), getFloorHandler, floorFilters),
		server.NewRouteWithFilters("PUT", floorPath(), updateFloorHandler, floorFilters),
		server.NewRouteWithFilters("DELETE", floorPath(), deleteFloorHandler, floorFilters),
	)
}

func loadFloorsAndFloorFromContext(kvStore store.Store, ctx server.RequestContext) error {
	err := loadBuildingAndBuildingFromContext(kvStore, ctx)
	if err != nil {
		return err
	}

	building, ok := ctx.UserValue(buildingUserKey).(gateway.Building)
	if !ok {
		return store.NotFound("unable to find building using context")
	}

	floors, err := kvStore.Floors(building)
	if err != nil {
		return err
	}

	id, ok := ctx.UserValue(floorID).(string)
	if !ok {
		return store.NotFound("unable to find floor key in the context")
	}

	floor, ok := floors[id]
	if !ok {
		return store.NotFound("unable to find floor using context")
	}

	ctx.SetUserValue(floorsUserKey, floors)
	ctx.SetUserValue(floorUserKey, floor)
	return nil
}

var findAndLoadFloor = func(kvStore store.Store, ctx server.RequestContext) error {
	err := loadFloorsAndFloorFromContext(kvStore, ctx)
	if err != nil {
		switch err.(type) {
		case store.NotFound:
			return notFound(ctx)
		default:
			return internalServerError(ctx, err)
		}
	}
	return ctx.Next()
}

var listFloorsHandler = func(store store.Store, ctx server.RequestContext) error {
	building, ok := ctx.UserValue(buildingUserKey).(gateway.Building)
	if !ok {
		return notFound(ctx)
	}

	floors, err := store.Floors(building)
	if err != nil {
		return internalServerError(ctx, err)
	}
	return ctx.JSONResponse(floors, http.StatusOK)
}

var createFloorHandler = func(store store.Store, ctx server.RequestContext) error {
	building, ok := ctx.UserValue(buildingUserKey).(gateway.Building)
	if !ok {
		return notFound(ctx)
	}

	floor, err := gateway.NewFloor(building, ctx.PostBody())
	if err != nil {
		return badRequest(ctx, err)
	}

	floors, err := store.Floors(building)
	if err != nil {
		return internalServerError(ctx, err)
	}

	floors[floor.ID()] = floor
	err = store.UpsertFloors(building, floors)
	if err != nil {
		return internalServerError(ctx, err)
	}
	return created(ctx)
}

var getFloorHandler = func(store store.Store, ctx server.RequestContext) error {
	floor, ok := ctx.UserValue(floorUserKey).(gateway.Floor)
	if !ok {
		return notFound(ctx)
	}

	return ctx.JSONResponse(floor, fasthttp.StatusOK)
}

var updateFloorHandler = func(store store.Store, ctx server.RequestContext) error {
	building, ok := ctx.UserValue(buildingUserKey).(gateway.Building)
	if !ok {
		return notFound(ctx)
	}

	floor, err := gateway.NewFloor(building, ctx.PostBody())
	if err != nil {
		return badRequest(ctx, err)
	}

	err = store.UpsertFloor(floor)
	if err != nil {
		return internalServerError(ctx, err)
	}
	return nil
}

var deleteFloorHandler = func(store store.Store, ctx server.RequestContext) error {
	floor, ok := ctx.UserValue(floorUserKey).(gateway.Floor)
	if !ok {
		return notFound(ctx)
	}

	err := store.DeleteFloor(floor)
	if err != nil {
		return internalServerError(ctx, err)
	}

	return nil
}
