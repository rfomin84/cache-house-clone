package public

import (
	"encoding/json"
	"fmt"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/valyala/fasthttp"
)

type SupplySidePlatformsController struct {
	SupplySidePlatformState *managers.SupplySidePlatformState
}

func (c *SupplySidePlatformsController) Index(ctx *fasthttp.RequestCtx) {

	ssps := c.SupplySidePlatformState.GetSupplySidePlatforms()
	sspResponse := managers.SupplySidePlatformsResponse{
		SupplySidePlatforms: ssps,
	}
	jsonResponse, err := json.Marshal(sspResponse)

	if err != nil {
		c.SupplySidePlatformState.Logger.Error(err.Error())
		ctx.Error("Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	_, _ = fmt.Fprintln(ctx, string(jsonResponse))
}
