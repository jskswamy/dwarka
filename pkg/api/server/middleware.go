package server

import (
	"fmt"
	"github.com/savsgio/atreugo/v10"
	"time"
)

const (
	startTime = "request.startTime"
)

type handlerInfo struct {
	method    []byte
	path      []byte
	startTime time.Time
}

var startMeasure = func(ctx *atreugo.RequestCtx) error {
	ctx.SetUserValue("info", handlerInfo{
		method:    ctx.Method(),
		path:      ctx.Path(),
		startTime: time.Now(),
	})
	return ctx.Next()
}

var stopMeasure = func(ctx *atreugo.RequestCtx) error {
	info, ok := ctx.UserValue("info").(handlerInfo)
	if ok {
		duration := time.Since(info.startTime)
		ctx.Logger().Printf("%s - %d - %s", info.path, ctx.Response.StatusCode(), timeUnit(duration))
	}
	return ctx.Next()
}

func timeUnit(duration time.Duration) string {
	us := duration.Microseconds()
	ms := duration.Milliseconds()
	ns := duration.Nanoseconds()
	if ns < 1000 {
		return fmt.Sprintf("%d ns", ns)
	} else if us < 1000 {
		return fmt.Sprintf("%d μs", us)
	}
	return fmt.Sprintf("%d ms", ms)
}
