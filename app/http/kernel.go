package http

import (
	"github.com/sunranlike/sunr/framework"
	"github.com/sunranlike/sunr/framework/gin"
)

// NewHttpEngine 创建了一个绑定了路由的Web引擎
func NewHttpEngine(container framework.Container) (*gin.Engine, error) {
	// 设置为Release，为的是默认在启动中不输出调试信息,默认是输出调试信息的
	gin.SetMode(gin.ReleaseMode)
	//// 默认启动一个Web引擎,这个方法调用New,并且传入了一个日志中间件logger和恢复中间件recoveery
	//r := gin.Default()
	r := gin.New() //New()返回一个默认的Engine,这个Engine嵌入了实现Container接口的结构体,
	// 这样我们的gin.Engine也可以使用container接口
	r.SetContainer(container)
	r.Use(gin.Logger(), gin.Recovery()) //使用两个全局中间件.

	// 业务绑定路由操作,绑定路由就是对一个uri绑定到handler,实际上不仅仅只是绑定
	//5.6我们更改了动态路由
	Routes(r)
	// 返回绑定路由后的Web引擎
	return r, nil
}
