package main

import (
	"coredemo/framework"
	"coredemo/framework/middleware"
	"net/http"
)

func main() {
	core := framework.NewCore()                                     //注册声明一个core,这个Newcore函数,core到底是什么?core是一个实现了ServeHttp方法的结构体,也就是说core实现了Handler接口
	core.Use(middleware.Recovery(), middleware.RecordRequsstTime()) //将中间件添加到slice中

	registerRouter(core)
	//注册core,目的是为了绑定foo与FooControllerHandler,让他们一个作为map一个作为value
	server := &http.Server{ //这个结构体的Handler字段就只去使用自己的Handler函数
		// 自定义的请求核心处理函数
		Handler: core,
		//Handler是一个http.Serve的接口,这个接口的方法集合只有一个方法:serveHTTP,会默认执行这个serveHTTP
		//如果你没有实现这个serveHTTP方法,就会使用默认的defaultServeHTTP方法,也就是我们的core需要一个有一个方法叫serveHTTP去被执行
		//我们也确实实现了这个接口
		// 请求监听地址
		Addr: ":8080",
	}
	//开始监听:这是个server的内置方法,主要调用net.Listen*("tcp", addr)
	//当然最底层调用的是TCPListrener
	//这里也说明了http是在tcp上面的方法,http是应用层,他通过tcp实现
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
