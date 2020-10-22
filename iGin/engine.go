package iGin

import (
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	routerMap map[string]HandlerFunc
}

//ServeHttp为单次请求的内容
func (engine *Engine) ServeHTTP(rspWriter http.ResponseWriter, req *http.Request) {
	targetUrl := req.Method + "_" + req.URL.Path
	if handler, ok := engine.routerMap[targetUrl]; ok {
		ctx := NewContext(rspWriter, req)
		handler(ctx)
	} else {
		http.Error(rspWriter, "404 not found", 404)
	}
}

func (engine *Engine) registerHandler(method string, url string, handler HandlerFunc) {
	targetStr := method + "_" + url
	engine.routerMap[targetStr] = handler
}

func (engine *Engine) Get(url string, handler HandlerFunc) {
	engine.registerHandler("GET", url, handler)
}

func (engine *Engine) Post(url string, handler HandlerFunc) {
	engine.registerHandler("POST", url, handler)
}

func (engine *Engine) Serve(port string) {
	_ = http.ListenAndServe(port, engine)
}

func (engine *Engine) New() {
	engine.routerMap = make(map[string]HandlerFunc)
}