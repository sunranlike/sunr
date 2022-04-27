package main

import (
	"context"
	"github.com/sunranlike/hade/framework/gin"
	"github.com/sunranlike/hade/framework/provider/demo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//core := framework.NewCore()                                     //注册声明一个core,这个Newcore函数,core到底是什么?core是一个实现了ServeHttp方法的结构体,也就是说core实现了Handler接口
	//core.Use(middleware.Recovery()) //将中间件添加到slice中
	//core.Use(middleware.CostTime())
	////为什么registerRouter可以实现core的路由功能?因为这个函数追究到底调用了了core的Get方法,这个方法可以把url注册到core的router中
	//registerRouter(core)
	core := gin.New()
	core.Bind(&demo.DemoServiceProvider{})
	//fmt.Println(core)
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

	//开启一个routine来启动服务,这样服务就不会再main routine中执行了，从而可以避免被杀死
	go func() {
		err := server.ListenAndServe() //让他在后台开个协程跑着
		if err != nil {
			return
		} //
	}()

	// //声明一个quit信号channel,会用来作为监听三个系统信号的通道
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	//如果调用者不执行执行那三个系统调用,main routine就会一直阻塞在这里
	<-quit //会阻塞在这里,知道有信号通知结束

	//写一个定时器，作为参数传入shutdown的话就会
	//timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//server.Shutdown 方法是个阻塞方法，一旦执行之后，
	//它会阻塞当前 Goroutine，并且在所有连接请求都结束之后(也就是server结束)，才继续往后执行
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

}
