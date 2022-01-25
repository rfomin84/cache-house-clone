package controllers

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type FeedsController struct {
}

func (c *FeedsController) Index(ctx *fasthttp.RequestCtx) {
	_, _ = fmt.Fprintln(ctx, "OK")
}
