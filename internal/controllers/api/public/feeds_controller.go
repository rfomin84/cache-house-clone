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

	var feedType managers.FeedType
	if len(queryParams["is_dsp"]) > 0 {
		value, err := strconv.Atoi(queryParams["is_dsp"][0])
		if err != nil {
			c.FeedState.Logger.Error(err.Error())
			ctx.Error("Error", fasthttp.StatusInternalServerError)
		}
		feedType = managers.FeedType(value)
	} else {
		feedType = managers.All
	}

	feeds := c.FeedState.GetFeeds(billingTypes, feedType)
	feedResponse := managers.FeedsResponse{
		Feeds: feeds,
	}
	jsonResponse, err := json.Marshal(feedResponse)

	if err != nil {
		c.FeedState.Logger.Error(err.Error())
		ctx.Error("Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	_, _ = fmt.Fprintln(ctx, string(jsonResponse))
}
