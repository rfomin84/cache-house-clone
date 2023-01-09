package public

import (
	"encoding/csv"
	"fmt"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type FeedsController struct {
	FeedState *managers.FeedState
}

func (c *FeedsController) Index(ctx echo.Context) error {

	queryParams := ctx.Request().URL.Query()

	billingTypes := queryParams["billing_types"]

	var feedType managers.FeedType
	if len(queryParams["is_dsp"]) > 0 {
		value, err := strconv.Atoi(queryParams["is_dsp"][0])
		if err != nil {
			c.FeedState.Logger.Error(err.Error())
			return ctx.String(http.StatusInternalServerError, "Error")
		}
		feedType = managers.FeedType(value)
	} else {
		feedType = managers.All
	}

	feeds := c.FeedState.GetFeeds(billingTypes, feedType)
	feedResponse := managers.FeedsResponse{
		Feeds: feeds,
	}

	return ctx.JSON(http.StatusOK, feedResponse)
}

func (c *FeedsController) FeedListTsv(ctx echo.Context) error {
	feeds := c.FeedState.GetFeedsNetworks()

	result := make([][]string, 0)
	result = append(result, []string{"campaign_id", "campaign_name"})
	for _, value := range feeds {
		row := make([]string, 0)
		row = append(row, strconv.Itoa(value.Id))
		row = append(row, fmt.Sprintf("%s (%d)", value.Name, value.Id))

		result = append(result, row)
	}

	ctx.Response().Header().Set(echo.HeaderContentType, "text/csv")
	ctx.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="feeds.tsv"`)

	writer := csv.NewWriter(ctx.Response().Writer)
	writer.Comma = '\t'
	_ = writer.WriteAll(result)

	return nil
}

func (c *FeedsController) ListAccountTsv(ctx echo.Context) error {
	feeds := c.FeedState.GetFeedsNetworks()

	result := make([][]string, 0)
	result = append(result, []string{"campaign_id", "account_id", "created_at"})
	for _, value := range feeds {
		row := make([]string, 0)
		row = append(row, strconv.Itoa(value.Id))
		row = append(row, strconv.Itoa(value.NetworkId))
		row = append(row, value.CreatedAt)

		result = append(result, row)
	}

	ctx.Response().Header().Set(echo.HeaderContentType, "text/csv")
	ctx.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="feed_account.tsv"`)

	writer := csv.NewWriter(ctx.Response().Writer)
	writer.Comma = '\t'
	_ = writer.WriteAll(result)

	return nil
}

func (c *FeedsController) ListNetworkTsv(ctx echo.Context) error {
	feeds := c.FeedState.GetFeedsNetworks()

	result := make([][]string, 0)
	result = append(result, []string{"campaign_id", "account_name"})
	for _, value := range feeds {
		row := make([]string, 0)
		row = append(row, strconv.Itoa(value.Id))
		row = append(row, value.NetworkName)

		result = append(result, row)
	}

	ctx.Response().Header().Set(echo.HeaderContentType, "text/csv")
	ctx.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="feed_network.tsv"`)

	writer := csv.NewWriter(ctx.Response().Writer)
	writer.Comma = '\t'
	_ = writer.WriteAll(result)

	return nil
}

func (c *FeedsController) FeedsAccountManagers(ctx echo.Context) error {
	feedsAccountManagers := c.FeedState.GetFeedsAccountManagers()

	result := make([][]string, 0)
	result = append(result, []string{"account_id", "campaign_id", "responsible_manager_id", "responsible_manager_name"})
	for _, value := range feedsAccountManagers {
		row := make([]string, 0)
		row = append(row, strconv.Itoa(value.AccountId))
		row = append(row, strconv.Itoa(value.CampaignId))
		managerId := ""
		if value.ResponsibleManagerId > 0 {
			managerId = strconv.Itoa(value.ResponsibleManagerId)
		}

		row = append(row, managerId)
		row = append(row, value.ResponsibleManagerName)

		result = append(result, row)
	}

	ctx.Response().Header().Set(echo.HeaderContentType, "text/csv")
	ctx.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="accounts_manager.tsv"`)

	writer := csv.NewWriter(ctx.Response().Writer)
	writer.Comma = '\t'
	_ = writer.WriteAll(result)

	return nil
}
