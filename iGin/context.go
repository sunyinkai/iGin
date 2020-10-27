package iGin

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	RawRsp http.ResponseWriter
	RawReq *http.Request
	//request
	Method string
	Path   string
	Params map[string]string
	//middleware
	index    int
	handlers []HandlerFunc
}

func NewContext(rspWriter http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		RawReq:   req,
		RawRsp:   rspWriter,
		Method:   req.Method,
		Path:     req.URL.Path,
		Params:   make(map[string]string),
		index:    -1,
		handlers: make([]HandlerFunc, 0),
	}
}

func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ { //因为不是所有的handler都会调用c.Next(),为了兼容
		c.handlers[c.index](c)
	}
}

func (c *Context) SetHeader(key, value string) {
	c.RawRsp.Header().Set(key, value)
}
func (c *Context) SetStatus(code int) {
	c.RawRsp.WriteHeader(code)
}

func (c *Context) PostForm(key string) string {
	return c.RawReq.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.RawReq.URL.Query().Get(key)
}

func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.RawRsp)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.RawRsp, err.Error(), 500)
	}
}

func (c *Context) Html(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	_, _ = c.RawRsp.Write([]byte(html))
}
