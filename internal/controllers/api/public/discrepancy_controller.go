package public

import (
	"encoding/json"
	"fmt"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/golang-module/carbon/v2"
	"github.com/valyala/fasthttp"
	"net/url"
	"strconv"
)

type DiscrepancyController struct {
	DiscrepancyState *managers.DiscrepancyState
}

func (c *DiscrepancyController) Index(ctx *fasthttp.RequestCtx) {
	queryParams, _ := url.ParseQuery(ctx.QueryArgs().String())

	billingTypes := queryParams["billing_types"]

	var feedType managers.FeedType
	if len(queryParams["is_dsp"]) > 0 {
		value, err := strconv.Atoi(queryParams["is_dsp"][0])
		if err != nil {
			c.DiscrepancyState.Logger.Error(err.Error())
			ctx.Error("Error", fasthttp.StatusInternalServerError)
		}
		feedType = managers.FeedType(value)
	} else {
		feedType = managers.All
	}

	startDate := carbon.Parse(string(ctx.QueryArgs().Peek("start_date")))
	endDate := carbon.Parse(string(ctx.QueryArgs().Peek("end_date")))

	result := c.DiscrepancyState.GetDiscrepancies(startDate.Carbon2Time(), endDate.Carbon2Time(), billingTypes, feedType)

	jsonResponse, err := json.Marshal(result)

	if err != nil {
		c.DiscrepancyState.Logger.Error(err.Error())
		ctx.Error("Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	_, _ = fmt.Fprintln(ctx, string(jsonResponse))
}
