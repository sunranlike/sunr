package framework

import (
	"context"
	"net/http"
)

// 自定义 Context
type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
}

//BaseContext()返回的就是标准库的ctx
func (ctx *Context) BaseContext() context.Context {
	//直接使用request的Ctx,导入这个标准包,我就不用再去写新的实现了
	return ctx.request.Context()
}

func (ctx *Context) Done() <-chan struct{} {
	//同样Done函数也用BaseContext()的Done(),就不用我们去实现了
	return ctx.BaseContext().Done()
}
