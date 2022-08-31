package public

import (
	"encoding/json"
	"fmt"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/valyala/fasthttp"
	"net/url"
	"strconv"
	"time"
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

	startDate, _ := time.Parse("2006-01-02", queryParams["start_date"][0])
	endDate, _ := time.Parse("2006-01-02", queryParams["end_date"][0])

	fmt.Println(startDate)
	fmt.Println(endDate)

	fmt.Println(billingTypes, feedType)

	result := c.DiscrepancyState.GetDiscrepancies(startDate, endDate, billingTypes, feedType)

	jsonResponse, err := json.Marshal(result)

	if err != nil {
		c.DiscrepancyState.Logger.Error(err.Error())
		ctx.Error("Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	_, _ = fmt.Fprintln(ctx, string(jsonResponse))
}
