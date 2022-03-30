package controllers

import (
	"fmt"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/valyala/fasthttp"
)

type HomeController struct {
	FeedState *managers.FeedState
}

func (c *HomeController) Index(ctx *fasthttp.RequestCtx) {
	if len(c.FeedState.Feeds) == 0 {
		ctx.Error("", fasthttp.StatusInternalServerError)
		return
	}
	fmt.Fprintln(ctx, "OK")
}
