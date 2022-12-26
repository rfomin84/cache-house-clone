package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
)

func AuthMiddleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		token := string(ctx.QueryArgs().Peek("api_token"))
		if token == os.Getenv("API_TOKEN") {
			handler(ctx)
			return
		}
		ctx.Error("Error", fasthttp.StatusUnauthorized)
		return
	})
}

//Auth middleware
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("api_token")

		if token == os.Getenv("API_TOKEN") {
			return next(c)
		}
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
	}
}
