package main

import (
	"github.com/sunranlike/hade/app/console"
	"github.com/sunranlike/hade/app/http"
	"github.com/sunranlike/hade/framework"
	"github.com/sunranlike/hade/framework/provider/app"
	"github.com/sunranlike/hade/framework/provider/config"
	"github.com/sunranlike/hade/framework/provider/env"
	"github.com/sunranlike/hade/framework/provider/id"
	"github.com/sunranlike/hade/framework/provider/kernel"
	"github.com/sunranlike/hade/framework/provider/log"
	"github.com/sunranlike/hade/framework/provider/trace"
)

func main() {

	//初始化服务容器:调用的是NewHadeContainer()方法，该方法返回一个HadeContainer接口的指针
	//这个hadeContainer实现了framework.Container接口的五个方法
	container := framework.NewHadeContainer()

	//绑定框架级服务提供者 (业务级别服务不在这里绑定,在更下层的函数绑定): 目录结构服务 HadeAppProvider
	//我们调用了Container接口的一个方法:Bind,将一个service provider和我们的容器做绑定.这个服务需要实现serviceProvider接口

	container.Bind(&app.HadeAppProvider{})       //appProvider
	container.Bind(&env.HadeEnvProvider{})       //环境变量服务
	container.Bind(&config.HadeConfigProvider{}) //配置读取服务
	container.Bind(&id.HadeIDProvider{})
	container.Bind(&trace.HadeTraceProvider{})
	container.Bind(&log.HadeLogServiceProvider{})
	//先初始化engine实例才可以传入HadeKernelP结构体内,因为新建http可能会失败,所以要handle err
	//之前并没有绑定这个本地分布式抢占系统，所以会提示 contract hade:distributed have not register
	//container.Bind(&distributed.LocalDistributedProvider{})
	if engine, err := http.NewHttpEngine(); err == nil {

		//这里又绑定一个框架级别服务?why,其实就是对上面的那个gin.engige进行绑定,上面使用一个common ok 语法接收,如果初始化成功
		//将这个web服务engine绑定到我们的HadeKernelProvider结构体,这个结构体是一个实现了ServiceProvider的实例框架级服务,其中有一个gin.Engine字段
		//NewHttpEngine返回的engine还对engine绑定了两个全局中间件,并且还配置类路由,
		//并且还绑定了一个业务级别服务demoService到engine中。
		//所以要区分好业务服务和框架服务，业务服务要在路由中绑定到gin.engine中，而框架级服务要绑定在主题container中
		//最终这个web服务又是一个框架级别服务,就可以绑定到我们的容器中,这个HadeKernelProvider框架级服务

		container.Bind(&kernel.HadeKernelProvider{HttpEngine: engine})

		//一定要理解container绑定的是框架级别服务,web服务是一个框架级别服务,绑定到里面,然后我们的demoservice是一个业务级别
		//服务,需要绑定的是web服务的主题gin.Engine的容器中,实现了业务和框架的分离
	}

	// 运行 root 命令
	//如果说上面是将服务绑定到容器之中,那么接下来就是对cobra的改造
	//执行该函数会声明一个根command,我们所有的命令行都是根据这个command
	console.RunCommand(container)

}
