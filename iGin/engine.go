package iGin

import (
	"log"
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	routerManager RouterManager
}

//ServeHttp为单次请求的内容
func (engine *Engine) ServeHTTP(rspWriter http.ResponseWriter, req *http.Request) {
	targetUrl := req.Method + "_" + req.URL.Path
	if ok, handlers := engine.routerManager.Query(targetUrl); ok {
		log.Printf("handlers:%+v\n", handlers)
		ctx := NewContext(rspWriter, req)
		ctx.handlers = handlers
		ctx.Next()
	} else {
		http.Error(rspWriter, "404 not found", 404)
	}
}

//注册试图函数
func (engine *Engine) registerViewFunc(method string, url string, handler HandlerFunc) {
	targetStr := method + "_" + url
	_, _ = engine.routerManager.InsertViewFunc(targetStr, handler)
}

//注册中间件
func (engine *Engine) registerMiddleWare(method string, url string, handlers []HandlerFunc) {
	targetStr := method + "_" + url
	_, _ = engine.routerManager.InsertMiddleWare(targetStr, handlers)
}

func (engine *Engine) Group(method string, url string, handlers ...HandlerFunc) {
	engine.registerMiddleWare(method, url, handlers)
}

func (engine *Engine) Get(url string, handler HandlerFunc) {
	engine.registerViewFunc("GET", url, handler)
}

func (engine *Engine) Post(url string, handler HandlerFunc) {
	engine.registerViewFunc("POST", url, handler)
}

func (engine *Engine) Serve(port string) {
	_ = http.ListenAndServe(port, engine)
}

func (engine *Engine) New() {
	//engine.routerMap = make(map[string]HandlerFunc)
}
