package main

import (
	"fmt"
	"iGin/iGin"
)

func main() {
	var engine iGin.Engine
	engine.New()
	engine.Get("/", func(ctx *iGin.Context) {
		result := fmt.Sprintf("{'code':'0','path':'%s','host':'%s'}", ctx.Path, ctx.RawReq.Host)
		ctx.Json(200, result)
	})
	engine.Get("/hello/", func(ctx *iGin.Context) {
		ctx.Html(200, "<h1>you are hello </h1>")
	})
	engine.Serve(":8080")
}
