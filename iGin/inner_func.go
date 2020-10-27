package iGin

import "net/http"

func DefaultNotFound(ctx *Context) {
	http.Error(ctx.RawRsp, "404 not found", 404)
}
