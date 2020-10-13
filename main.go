package main

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(rspWriter http.ResponseWriter, req *http.Request);

type Engine struct {
	routerMap map[string]HandlerFunc
}

func (engine *Engine) registerHandler(method string, url string, handler HandlerFunc) {
	targetStr := method + "_" + url
	engine.routerMap[targetStr] = handler
}

//ServeHttp为单次请求的操作
func (engine *Engine) ServeHTTP(rspWriter http.ResponseWriter, req *http.Request) {
	targetUrl := req.Method + "_" + req.URL.Path
	if hander, ok := engine.routerMap[targetUrl]; ok {
		hander(rspWriter, req)
	} else {
		fmt.Fprintf(rspWriter, "404 not found\n")
	}
}

func (engine *Engine) Get(url string, handler HandlerFunc) {
	engine.registerHandler("GET", url, handler)
}

func (engine *Engine) Post(url string, handler HandlerFunc) {
	engine.registerHandler("POST", url, handler)
}

func (engine *Engine) Serve(port string) {
	http.ListenAndServe(port, engine)
}

func (engine *Engine) New() {
	engine.routerMap = make(map[string]HandlerFunc)
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		answer := "hello,you are in index"
		_, _ = writer.Write([]byte(answer))
	})
	var engine Engine
	engine.New()
	engine.Get("/", func(rspWriter http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rspWriter, "hello,/")
	})
	engine.Get("/iGin", func(rspWriter http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rspWriter, "hello,iGin")
	})
	engine.Serve(":8080")
}
