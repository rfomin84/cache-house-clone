package public

import (
	"encoding/csv"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type NetworkController struct {
	NetworkState *managers.NetworkState
}

func (c *NetworkController) Index(ctx echo.Context) error {
	networks := c.NetworkState.GetNetworks()
	networkResponse := managers.NetworksResponse{
		Networks: networks,
	}

	return ctx.JSON(http.StatusOK, networkResponse)
}

func (c *NetworkController) Tsv(ctx echo.Context) error {
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

	ctx.Response().Header().Set(echo.HeaderContentType, "text/csv")
	ctx.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="feeds.tsv"`)

	writer := csv.NewWriter(ctx.Response().Writer)
	writer.Comma = '\t'
	_ = writer.WriteAll(result)

	return nil
}
