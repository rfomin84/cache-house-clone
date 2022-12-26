package controllers

import (
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/labstack/echo/v4"
	"net/http"
)

type HomeController struct {
	FeedState    *managers.FeedState
	SspState     *managers.SupplySidePlatformState
	DiscrepState *managers.DiscrepancyState
}

func (c *HomeController) Index(ctx echo.Context) error {

	countFeeds := len(c.FeedState.Feeds)
	countSsp := len(c.SspState.SupplySidePlatforms)
	countDiscrep := len(c.DiscrepState.Discrepancies)

	if countFeeds > 0 && countSsp > 0 && countDiscrep > 0 {
		return ctx.String(http.StatusOK, "OK")
	} else {
		return ctx.String(http.StatusInternalServerError, "")
	}
}
