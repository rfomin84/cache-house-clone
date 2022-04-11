package middleware

import (
	"github.com/valyala/fasthttp"
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
