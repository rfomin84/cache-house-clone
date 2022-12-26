package public

import (
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"strconv"
)

type DiscrepancyController struct {
	DiscrepancyState *managers.DiscrepancyState
}

func (c *DiscrepancyController) Index(ctx echo.Context) error {
	queryParams, _ := url.ParseQuery(ctx.Request().URL.String())

	billingTypes := queryParams["billing_types"]

	var feedType managers.FeedType
	if len(queryParams["is_dsp"]) > 0 {
		value, err := strconv.Atoi(queryParams["is_dsp"][0])
		if err != nil {
			c.DiscrepancyState.Logger.Error(err.Error())
			return ctx.String(http.StatusInternalServerError, "Error")
		}
		feedType = managers.FeedType(value)
	} else {
		feedType = managers.All
	}

	startDate := carbon.Parse(ctx.QueryParam("start_date"))
	endDate := carbon.Parse(ctx.QueryParam("end_date"))

	result := c.DiscrepancyState.GetDiscrepancies(startDate.Carbon2Time(), endDate.Carbon2Time(), billingTypes, feedType)

	ctx.Response().Header().Set(echo.HeaderContentType, "application/json")
	return ctx.JSON(http.StatusOK, result)
}
