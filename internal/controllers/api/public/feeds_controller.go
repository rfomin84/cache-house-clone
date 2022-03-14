package public

import (
	"encoding/json"
	"fmt"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/valyala/fasthttp"
	"net/url"
	"strconv"
)

type FeedsController struct {
	FeedState *managers.FeedState
}

func (c *FeedsController) Index(ctx *fasthttp.RequestCtx) {

	queryParams, _ := url.ParseQuery(ctx.QueryArgs().String())

	billingTypes := queryParams["billing_types"]

	isDsp := false
	if len(queryParams["is_dsp"]) > 0 {
		isDsp, _ = strconv.ParseBool(queryParams["is_dsp"][0])
	}

	feeds := c.FeedState.GetFeeds(billingTypes, isDsp)
	jsonResponse, err := json.Marshal(feeds)

	if err != nil {
		c.FeedState.Logger.Error(err.Error())
		ctx.Error("Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	_, _ = fmt.Fprintln(ctx, string(jsonResponse))
}
