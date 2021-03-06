package main

import (
	"fmt"
	"iGin/iGin"
	"log"
	"time"
)

func init() {
	log.SetPrefix("main ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	engine := iGin.NewIGinEngine()

	loggerMiddle := func() iGin.HandlerFunc {
		return func(ctx *iGin.Context) {
			t := time.Now()
			ctx.Next()
			log.Printf("response %s in %v", ctx.Path, time.Since(t))
		}
	}
	testMiddle1 := func() iGin.HandlerFunc {
		return func(ctx *iGin.Context) {
			log.Printf("A1")
			ctx.Next()
			log.Printf("A2")
		}
	}
	testMiddle2 := func() iGin.HandlerFunc {
		return func(ctx *iGin.Context) {
			log.Printf("B1")
			ctx.Next()
			log.Printf("B2")
		}
	}
	testMiddle3 := func() iGin.HandlerFunc {
		return func(ctx *iGin.Context) {
			log.Printf("C1")
			ctx.Next()
			log.Printf("C2")
		}
	}

	engine.Use("GET", "/", loggerMiddle())
	engine.Use("GET", "/", testMiddle1(), testMiddle2(), testMiddle3())

	engine.Get("/", func(ctx *iGin.Context) {
		result := fmt.Sprintf("{'code':'0','path':'%s','host':'%s'}", ctx.Path, ctx.RawReq.Host)
		ctx.Json(200, result)
	})
	engine.Get("/hello", func(ctx *iGin.Context) {
		log.Printf("log,headers:%+v", ctx.RawRsp.Header())
		ctx.Redirect(301, "/world")
		log.Printf("log,headers:%+v", ctx.RawRsp.Header())
		htmlStr := fmt.Sprintf("<h1>you are now in hello name=%s</h1>", ctx.Query("name"))
		ctx.Html(200, htmlStr)
	})
	engine.Get("/world/", func(ctx *iGin.Context) {
		htmlStr := fmt.Sprintf("<h1>you are now in world,name=%s </h1>", ctx.Query("name"))
		ctx.Redirect(301, "https://baidu.com")
		ctx.Html(200, htmlStr)
	})
	engine.Get("/world/a", func(ctx *iGin.Context) {
		htmlStr := fmt.Sprintf("<h1>you are now in world a,name=%s </h1>", ctx.Query("name"))
		ctx.Html(200, htmlStr)
	})
	engine.Get("/tt/:name/:sex", func(ctx *iGin.Context) {
		ctx.Json(200, ctx.Params)
	})
	engine.Get("/static/*filepath", func(ctx *iGin.Context) {
		ctx.Json(200, ctx.Params)
	})
	engine.Serve("0.0.0.0:8080")
}
