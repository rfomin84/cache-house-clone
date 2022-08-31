package public

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/valyala/fasthttp"
	"strconv"
)

type NetworkController struct {
	NetworkState *managers.NetworkState
}

func (c *NetworkController) Index(ctx *fasthttp.RequestCtx) {
	networks := c.NetworkState.GetNetworks()
	networkResponse := managers.NetworksResponse{
		Networks: networks,
	}

	jsonResponse, err := json.Marshal(networkResponse)

	if err != nil {
		c.NetworkState.Logger.Error(err.Error())
		ctx.Error("Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	_, _ = fmt.Fprintln(ctx, string(jsonResponse))
}

func (c *NetworkController) Tsv(ctx *fasthttp.RequestCtx) {
	networks := c.NetworkState.GetNetworks()
	networkResponse := managers.NetworksResponse{
		Networks: networks,
	}
	result := make([][]string, 0)
	result = append(result, []string{"account_id", "name", "created_at"})
	for _, value := range networkResponse.Networks {
		row := make([]string, 0)
		row = append(row, strconv.Itoa(value.AccountId))
		row = append(row, value.Timezone)
		row = append(row, value.FeedIds)

		result = append(result, row)
	}

	writer := csv.NewWriter(ctx.Response.BodyWriter())
	writer.Comma = '\t'
	_ = writer.WriteAll(result)

	ctx.Response.Header.Set("Content-Type", "text/csv")
	ctx.Response.Header.Set("Content-Disposition", `attachment; filename="networks.tsv"`)
}
