package api

import (
	"fmt"
	"net/http"
	"path"

	"github.com/valyala/fasthttp"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/api/server"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/api/view"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
)

const (
	roomID       = "room-id"
	roomUserKey  = "room"
	roomsUserKey = "rooms"
)

func roomPath() string {
	return path.Join(roomsBasePath(), fmt.Sprintf(":%s", roomID))
}

func roomsBasePath() string {
	return path.Join(floorPath(), "rooms")
}

func init() {
	roomFilters := &server.Filters{Before: []server.ResponseHandler{findAndLoadRoom}}
	floorFilters := &server.Filters{Before: []server.ResponseHandler{findAndLoadFloor}}
	AddRoute(
		server.NewRouteWithFilters("GET", roomsBasePath(), listRoomsHandler, floorFilters),
		server.NewRouteWithFilters("POST", roomsBasePath(), createRoomHandler, floorFilters),
		server.NewRouteWithFilters("GET", roomPath(), getRoomHandler, roomFilters),
		server.NewRouteWithFilters("PUT", roomPath(), updateRoomHandler, roomFilters),
		server.NewRouteWithFilters("DELETE", roomPath(), deleteRoomHandler, roomFilters),
	)
}

func loadRoomsAndRoomFromContext(kvStore store.Store, ctx server.RequestContext) error {
	err := loadFloorsAndFloorFromContext(kvStore, ctx)
	if err != nil {
		return nil
	}

	floor, ok := ctx.UserValue(floorUserKey).(gateway.Floor)
	if !ok {
		return store.NotFound("unable to find floor using context")
	}

	rooms, err := kvStore.Rooms(floor)
	if err != nil {
		return err
	}

	id, ok := ctx.UserValue(roomID).(string)
	if !ok {
		return store.NotFound("unable to find room key in the context")
	}

	room, ok := rooms[id]
	if !ok {
		return store.NotFound("unable to find room using context")
	}

	ctx.SetUserValue(roomsUserKey, rooms)
	ctx.SetUserValue(roomUserKey, room)
	return nil
}

var findAndLoadRoom = func(kvStore store.Store, ctx server.RequestContext) error {
	err := loadRoomsAndRoomFromContext(kvStore, ctx)
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

var listRoomsHandler = func(store store.Store, ctx server.RequestContext) error {
	floor, ok := ctx.UserValue(floorUserKey).(gateway.Floor)
	if !ok {
		return notFound(ctx)
	}

	rooms, err := store.Rooms(floor)
	if err != nil {
		return internalServerError(ctx, err)
	}
	return ctx.JSONResponse(view.NewRooms(rooms), http.StatusOK)
}

var createRoomHandler = func(store store.Store, ctx server.RequestContext) error {
	floor, ok := ctx.UserValue(floorUserKey).(gateway.Floor)
	if !ok {
		return notFound(ctx)
	}

	room, err := view.Convert(floor, ctx.PostBody())
	if err != nil {
		return badRequest(ctx, err)
	}

	rooms, err := store.Rooms(floor)
	if err != nil {
		return internalServerError(ctx, err)
	}

	rooms[room.ID()] = room
	err = store.UpsertRooms(floor, rooms)
	if err != nil {
		return internalServerError(ctx, err)
	}
	return created(ctx, room.ID())
}

var getRoomHandler = func(store store.Store, ctx server.RequestContext) error {
	room, ok := ctx.UserValue(roomUserKey).(gateway.Room)
	if !ok {
		return notFound(ctx)
	}

	return ctx.JSONResponse(view.NewRoom(room), fasthttp.StatusOK)
}

var updateRoomHandler = func(store store.Store, ctx server.RequestContext) error {
	floor, ok := ctx.UserValue(floorUserKey).(gateway.Floor)
	if !ok {
		return notFound(ctx)
	}

	room, err := view.Convert(floor, ctx.PostBody())
	if err != nil {
		return badRequest(ctx, err)
	}

	err = store.UpsertRoom(room)
	if err != nil {
		return internalServerError(ctx, err)
	}

	return nil
}

var deleteRoomHandler = func(store store.Store, ctx server.RequestContext) error {
	room, ok := ctx.UserValue(roomUserKey).(gateway.Room)
	if !ok {
		return notFound(ctx)
	}

	err := store.DeleteRoom(room)
	if err != nil {
		return internalServerError(ctx, err)
	}

	return nil
}
