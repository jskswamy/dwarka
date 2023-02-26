package api

import (
	"fmt"
	"net/http"
	"path"

	"github.com/valyala/fasthttp"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/api/server"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
)

const (
	buildingsBasePath = "/buildings"
	buildingID        = "building-id"
	buildingUserKey   = "building"
	buildingsUserKey  = "buildings"
)

func buildingPath() string {
	return path.Join(buildingsBasePath, fmt.Sprintf("{%s}", buildingID))
}

func init() {
	buildingFilters := &server.Filters{Before: []server.ResponseHandler{findAndLoadBuilding}}
	AddRoute(
		server.NewRoute("GET", buildingsBasePath, listBuildingsHandler),
		server.NewRoute("POST", buildingsBasePath, createBuildingHandler),
		server.NewRouteWithFilters("GET", buildingPath(), getBuildingHandler, buildingFilters),
		server.NewRouteWithFilters("PUT", buildingPath(), updateBuildingHandler, buildingFilters),
		server.NewRouteWithFilters("DELETE", buildingPath(), deleteBuildingHandler, buildingFilters),
	)
}

func loadBuildingAndBuildingFromContext(kvStore store.Store, ctx server.RequestContext) error {
	buildings, err := kvStore.Buildings()
	if err != nil {
		return err
	}

	id, ok := ctx.UserValue(buildingID).(string)
	if !ok {
		return store.NotFound("unable to find building key in the context")
	}

	building, ok := buildings[id]
	if !ok {
		return store.NotFound("unable to find building in store")
	}

	ctx.SetUserValue(buildingsUserKey, buildings)
	ctx.SetUserValue(buildingUserKey, building)
	return nil
}

var findAndLoadBuilding = func(kvStore store.Store, ctx server.RequestContext) error {
	err := loadBuildingAndBuildingFromContext(kvStore, ctx)
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

var listBuildingsHandler = func(store store.Store, ctx server.RequestContext) error {
	buildings, err := store.Buildings()
	if err != nil {
		return internalServerError(ctx, err)
	}
	return ctx.JSONResponse(buildings, http.StatusOK)
}

var createBuildingHandler = func(store store.Store, ctx server.RequestContext) error {
	building, err := gateway.NewBuilding(ctx.PostBody())
	if err != nil {
		return badRequest(ctx, err)
	}

	buildings, err := store.Buildings()
	if err != nil {
		return internalServerError(ctx, err)
	}

	buildings[building.ID()] = building
	err = store.UpsertBuildings(buildings)
	if err != nil {
		return internalServerError(ctx, err)
	}
	return created(ctx, building.ID())
}

var getBuildingHandler = func(store store.Store, ctx server.RequestContext) error {
	building, ok := ctx.UserValue(buildingUserKey).(gateway.Building)
	if !ok {
		return notFound(ctx)
	}

	return ctx.JSONResponse(building, fasthttp.StatusOK)
}

var updateBuildingHandler = func(store store.Store, ctx server.RequestContext) error {
	building, err := gateway.NewBuilding(ctx.PostBody())
	if err != nil {
		return badRequest(ctx, err)
	}

	err = store.UpsertBuilding(building)
	if err != nil {
		return internalServerError(ctx, err)
	}

	return nil
}

var deleteBuildingHandler = func(store store.Store, ctx server.RequestContext) error {
	building, ok := ctx.UserValue(buildingUserKey).(gateway.Building)
	if !ok {
		return notFound(ctx)
	}

	err := store.DeleteBuilding(building)
	if err != nil {
		return internalServerError(ctx, err)
	}

	return nil
}
