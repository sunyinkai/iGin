package main

import (
	"fmt"
	"iGin/iGin"
	"log"
	"time"
)

func main() {
	var engine iGin.Engine
	engine.New()
	engine.Group("GET", "/", func(ctx *iGin.Context) {
		t := time.Now()
		ctx.Next()
		log.Printf("response %s in %v", ctx.Path, time.Since(t))
	})
	engine.Group("GET", "/",
		func(ctx *iGin.Context) {
			log.Printf("A1")
			ctx.Next()
			log.Printf("A2")
		},
		func(ctx *iGin.Context) {
			log.Printf("B1")
			ctx.Next()
			log.Printf("B2")
		},
		func(ctx *iGin.Context) {
			log.Printf("C1")
			ctx.Next()
			log.Printf("C2")
		})

	engine.Get("/", func(ctx *iGin.Context) {
		result := fmt.Sprintf("{'code':'0','path':'%s','host':'%s'}", ctx.Path, ctx.RawReq.Host)
		ctx.Json(200, result)
	})
	engine.Get("/hello/", func(ctx *iGin.Context) {
		ctx.Html(200, "<h1>you are hello </h1>")
	})
	engine.Serve(":8080")
}
