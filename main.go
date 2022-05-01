package main

import (
	"github.com/sunranlike/hade/app/console"
	"github.com/sunranlike/hade/app/http"
	"github.com/sunranlike/hade/framework"
	"github.com/sunranlike/hade/framework/provider/app"
	"github.com/sunranlike/hade/framework/provider/kernel"
)

func main() {
	//调用gin的New函数,他返回的是一个Engine结构体
	//我们队Gin的Engine结构进行了修改,多了一个字段:container framework.Container
	//framework.Container是个结构,所以New返回的机构体必须实现这个接口,
	//因为Gin的New函数本来就实现了除了我们添加的字段的以外的初始化,
	//所以我们只需要在New函数中增加一个实现Container接口的结构体就是合法赋值
	//我们使用了NewHadeContainer()这个函数,他返回了一个实现Container接口的结构体
	//core := gin.New()

	//初始化服务容器:调用的是NewHadeContainer()方法
	container := framework.NewHadeContainer()

	//绑定服务提供者 : 目录结构服务 HadeAppProvider
	container.Bind(&app.HadeAppProvider{})
	if engine, err := http.NewHttpEngine(); err == nil { //先初始化engine实例才可以传入HadeKernelP结构体内
		container.Bind(&kernel.HadeKernelProvider{HttpEngine: engine})
	}

	// 运行 root 命令
	console.RunCommand(container)

	//core.Bind(&app.HadeAppProvider{BaseFolder: "/tmp"})

	//core.Bind(&demo.DemoServiceProvider{})
	//这里使用了两个全局中间件：
	//golang的net/http设计的一大特点就是特别容易构建中间件。
	//gin也提供了类似的中间件。需要注意的是中间件只对注册过的路由函数起作用。
	//对于分组路由，嵌套使用中间件，可以限定中间件的作用范围。
	//中间件分为全局中间件，单个路由中间件和群组中间件。

	//core.Use(gin.Recovery())
	//core.Use(middleware.Cost()) //这两个是全局中间件
	//
	//registerRouter(core)
	////fmt.Println(core)
	////注册core,目的是为了绑定foo与FooControllerHandler,让他们一个作为map一个作为value
	//server := &http.Server{ //这个结构体的Handler字段就只去使用自己的Handler函数
	//	// 自定义的请求核心处理函数
	//	Handler: core,
	//	//Handler是一个http.Serve的接口,这个接口的方法集合只有一个方法:serveHTTP,会默认执行这个serveHTTP
	//	//如果你没有实现这个serveHTTP方法,就会使用默认的defaultServeHTTP方法,也就是我们的core需要一个有一个方法叫serveHTTP去被执行
	//	//我们也确实实现了这个接口
	//	// 请求监听地址
	//	Addr: ":8080",
	//}
	//
	////开始监听:这是个server的内置方法,主要调用net.Listen*("tcp", addr)
	////当然最底层调用的是TCPListrener
	////这里也说明了http是在tcp上面的方法,http是应用层,他通过tcp实现
	//
	////开启一个routine来启动服务,这样服务就不会再main routine中执行了，从而可以避免被杀死
	//go func() {
	//	err := server.ListenAndServe() //让他在后台开个协程跑着
	//	if err != nil {
	//		return
	//	} //
	//}()
	//
	//// //声明一个quit信号channel,会用来作为监听三个系统信号的通道
	//quit := make(chan os.Signal)
	//// 监控信号：SIGINT, SIGTERM, SIGQUIT
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//
	////如果调用者不执行执行那三个系统调用,main routine就会一直阻塞在这里
	//<-quit //会阻塞在这里,知道有信号通知结束
	//
	////写一个定时器，作为参数传入shutdown的话就会
	////timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	////defer cancel()
	////server.Shutdown 方法是个阻塞方法，一旦执行之后，
	////它会阻塞当前 Goroutine，并且在所有连接请求都结束之后(也就是server结束)，才继续往后执行
	//if err := server.Shutdown(context.Background()); err != nil {
	//	log.Fatal("Server Shutdown:", err)
	//}

}
