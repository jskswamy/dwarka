package server

import (
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
)

//go:generate $PWD/scripts/mockgen $PWD/pkg/api/server/types.go $PWD/pkg/internal/mocks/api/server/types.go mockServer

// RequestContext represents necessary methods to handle HTTP request
type RequestContext interface {
	JSONResponse(body interface{}, statusCode ...int) error
	PostBody() []byte
	UserValue(key interface{}) interface{}
	SetStatusCode(statusCode int)
	SetBodyString(body string)
	Next() error
	SetUserValue(key interface{}, value interface{})
}

// ResponseHandler represents a function for responding to http request
type ResponseHandler func(store store.Store, ctx RequestContext) error

// RequestHandler represents a function for handling http request
type RequestHandler func(httpMethod, url string, handler ResponseHandler)

// Route represents a http route represented by httpMethod, url and
// response handler
type Route struct {
	httpMethod string
	url        string
	handler    ResponseHandler
	filters    *Filters
}

// Filters like middlewares, but for specific paths.
// It will be executed before and after the view defined in the path
// in addition of the general middlewares
type Filters struct {
	Before []ResponseHandler
	After  []ResponseHandler
}

// NewRoute returns a Route from httpMethod, url and handler
func NewRoute(httpMethod, url string, handler ResponseHandler) Route {
	return Route{
		httpMethod: httpMethod,
		url:        url,
		handler:    handler,
	}
}

// NewRouteWithFilters returns a Route from httpMethod, url and handler
func NewRouteWithFilters(httpMethod, url string, handler ResponseHandler, filters *Filters) Route {
	return Route{
		httpMethod: httpMethod,
		url:        url,
		handler:    handler,
		filters:    filters,
	}
}
