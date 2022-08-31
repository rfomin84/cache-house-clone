package public

import (
	"encoding/csv"
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

func (c *FeedsController) FeedListTsv(ctx *fasthttp.RequestCtx) {
	feeds := c.FeedState.GetFeeds(nil, managers.All)

	result := make([][]string, 0)
	result = append(result, []string{"campaign_id", "campaign_name"})
	for _, value := range feeds {
		row := make([]string, 0)
		row = append(row, strconv.Itoa(value.Id))
		row = append(row, fmt.Sprintf("%s (%d)", value.Name, value.Id))

		result = append(result, row)
	}

	writer := csv.NewWriter(ctx.Response.BodyWriter())
	writer.Comma = '\t'
	_ = writer.WriteAll(result)
}

func (c *FeedsController) ListAccountTsv(ctx *fasthttp.RequestCtx) {
	feeds := c.FeedState.GetFeeds(nil, managers.All)

	result := make([][]string, 0)
	result = append(result, []string{"campaign_id", "account_id", "created_at"})
	for _, value := range feeds {
		row := make([]string, 0)
		row = append(row, strconv.Itoa(value.Id))
		row = append(row, strconv.Itoa(value.AccountId))
		row = append(row, value.CreatedAt)

		result = append(result, row)
	}

	writer := csv.NewWriter(ctx.Response.BodyWriter())
	writer.Comma = '\t'
	_ = writer.WriteAll(result)

	ctx.Response.Header.Set("Content-Type", "text/csv")
	ctx.Response.Header.Set("Content-Disposition", `attachment; filename="feed_account.tsv"`)
}

func (c *FeedsController) ListNetworkTsv(ctx *fasthttp.RequestCtx) {
	feeds := c.FeedState.GetFeeds(nil, managers.All)

	result := make([][]string, 0)
	result = append(result, []string{"campaign_id", "account_name"})
	for _, value := range feeds {
		row := make([]string, 0)
		row = append(row, strconv.Itoa(value.Id))
		row = append(row, value.AccountName)

		result = append(result, row)
	}

	writer := csv.NewWriter(ctx.Response.BodyWriter())
	writer.Comma = '\t'
	_ = writer.WriteAll(result)

	ctx.Response.Header.Set("Content-Type", "text/csv")
	ctx.Response.Header.Set("Content-Disposition", `attachment; filename="feed_network.tsv"`)
}
