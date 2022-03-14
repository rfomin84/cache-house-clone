package controllers

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type HomeController struct {
}

func (c *HomeController) Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprintln(ctx, "OK")
}
