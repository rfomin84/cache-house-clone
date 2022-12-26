package public

import (
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SupplySidePlatformsController struct {
	SupplySidePlatformState *managers.SupplySidePlatformState
}

func (c *SupplySidePlatformsController) Index(ctx echo.Context) error {

	ssps := c.SupplySidePlatformState.GetSupplySidePlatforms()
	sspResponse := managers.SupplySidePlatformsResponse{
		SupplySidePlatforms: ssps,
	}

	return ctx.JSON(http.StatusOK, sspResponse)
}
