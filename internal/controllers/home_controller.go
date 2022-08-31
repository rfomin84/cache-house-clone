package controllers

import (
	"fmt"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/valyala/fasthttp"
)

type HomeController struct {
	FeedState    *managers.FeedState
	SspState     *managers.SupplySidePlatformState
	DiscrepState *managers.DiscrepancyState
}

func (c *HomeController) Index(ctx *fasthttp.RequestCtx) {

	countFeeds := len(c.FeedState.Feeds)
	countSsp := len(c.SspState.SupplySidePlatforms)
	countDiscrep := len(c.DiscrepState.Discrepancies)

	if countFeeds > 0 && countSsp > 0 && countDiscrep > 0 {
		fmt.Fprintln(ctx, "OK")
	} else {
		ctx.Error("", fasthttp.StatusInternalServerError)
		return
	}
}
