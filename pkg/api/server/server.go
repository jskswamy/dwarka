package server

import (
	"fmt"
	"github.com/savsgio/atreugo/v11"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
	"net"
	"time"
)

//go:generate $PWD/scripts/mockgen $PWD/pkg/api/server/server.go $PWD/pkg/internal/mocks/api/server/server.go mockServer

// Server represents necessary method to serve
// and handle HTTP traffic
type Server interface {
	ListenAndServe() error
	Serve(ln net.Listener) error
	Path(route Route)
}

// HTTPServer represents atreugo server backed by libkv/PersistentStore
type HTTPServer struct {
	atreugo *atreugo.Atreugo
	store   store.Store
}

// ListenAndServe binds the http server to the bind-address/bind-port
// and startMeasure serving http traffic
func (server HTTPServer) ListenAndServe() error {
	return server.atreugo.ListenAndServe()
}

// Serve serves incoming connections from the given listener.
//
// Serve blocks until the given listener returns permanent error.
//
// If we use a custom Listener, will be updated your atreugo configuration
// with the Listener address automatically
func (server HTTPServer) Serve(ln net.Listener) error {
	return server.atreugo.Serve(ln)
}

func (server HTTPServer) filters(handlers []ResponseHandler) []atreugo.Middleware {
	filters := make([]atreugo.Middleware, 0, len(handlers))
	for _, filter := range handlers {
		filters = append(filters, func(ctx *atreugo.RequestCtx) error {
			return filter(server.store, ctx)
		})
	}
	return filters
}

// Path binds a route to HTTPServer for handling request
func (server HTTPServer) Path(route Route) {
	path := server.atreugo.Path(route.httpMethod, route.url, func(ctx *atreugo.RequestCtx) error {
		return route.handler(server.store, ctx)
	})

	if route.filters != nil {
		path.Middlewares(atreugo.Middlewares{
			Before: server.filters(route.filters.Before),
			After:  server.filters(route.filters.After),
		})
	}
}

// NewHTTPServer returns a abstracted HTTP server
func NewHTTPServer(host, port string, store store.Store) Server {
	config := atreugo.Config{
		Name:              "dwarka",
		ReduceMemoryUsage: false,
		Addr:              fmt.Sprintf("%s:%s", host, port),
		GracefulShutdown:  true,
		WriteTimeout:      time.Second * 1,
		ReadTimeout:       time.Second * 1,
	}
	server := atreugo.New(config)
	server.UseBefore(startMeasure)
	server.UseAfter(stopMeasure)
	return &HTTPServer{atreugo: server, store: store}
}
