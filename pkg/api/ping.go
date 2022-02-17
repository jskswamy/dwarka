package api

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/api/server"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
)

func init() {
	AddRoute(server.NewRoute("GET", "/ping", pingHandler))
}

var pingHandler = func(store store.Store, ctx server.RequestContext) error {
	since, err := store.Uptime()
	if err != nil {
		status := map[string]string{"error": fmt.Sprintf("there was problem when reading value from store, reason: %v", err)}
		return ctx.JSONResponse(status, fasthttp.StatusInternalServerError)
	}

	status := map[string]interface{}{"status": since}
	return ctx.JSONResponse(status, fasthttp.StatusOK)
}
