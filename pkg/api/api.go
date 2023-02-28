package api

import (
	"github.com/valyala/fasthttp"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/api/server"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
)

var routes []server.Route

// AddRoute add route to list of known routes
func AddRoute(route ...server.Route) {
	routes = append(routes, route...)
}

// NewServer returns an abstracted http server with all
// the known routes added
func NewServer(host, port string, store store.Store) server.Server {
	httpServer := server.NewHTTPServer(host, port, store)
	for _, route := range routes {
		httpServer.Path(route)
	}
	return httpServer
}

func internalServerError(ctx server.RequestContext, err error) error {
	return ctx.JSONResponse(map[string]string{"error": err.Error()}, fasthttp.StatusInternalServerError)
}

func badRequest(ctx server.RequestContext, err error) error {
	return ctx.JSONResponse(map[string]string{"error": err.Error()}, fasthttp.StatusBadRequest)
}

func created(ctx server.RequestContext, id string) error {
	return ctx.JSONResponse(map[string]string{"id": id}, fasthttp.StatusCreated)
}

func notFound(ctx server.RequestContext) error {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	return nil
}

func conflict(ctx server.RequestContext) error {
	ctx.SetStatusCode(fasthttp.StatusConflict)
	return nil
}
