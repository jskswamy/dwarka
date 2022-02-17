package server

import (
	"fmt"
	"github.com/savsgio/atreugo/v10"
	"time"
)

const (
	startTime = "request.startTime"
)

var startMeasure = func(ctx *atreugo.RequestCtx) error {
	ctx.SetUserValue(startTime, time.Now())
	return ctx.Next()
}

var stopMeasure = func(ctx *atreugo.RequestCtx) error {
	startTime := ctx.UserValue(startTime).(time.Time)
	duration := time.Since(startTime)
	ctx.Logger().Printf("%d - %s", ctx.Response.StatusCode(), timeUnit(duration))
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
