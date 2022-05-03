package command

import (
	"context"
	"github.com/sunranlike/hade/framework/cobra"
	"github.com/sunranlike/hade/framework/contract"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// AppCommand 是命令行参数第一级为 app 的命令，它没有实际功能，只是打印帮助文档

// initAppCommand 初始化app命令和其子命令
func initAppCommand() *cobra.Command {
	appCommand.AddCommand(appStartCommand) //有一个父子关系,父亲是appCommand.子是appStartCommand
	//对应的父命令式app,子命令是start
	return appCommand
}

// AppCommand 是命令行参数第一级为app的命令，它没有实际功能，只是打印帮助文档,关键在于子命令
var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(c *cobra.Command, args []string) error {
		// 打印帮助文档
		c.Help()
		return nil
	},
}

// appStartCommand 启动一个Web服务,web服务是框架级服务.
var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个Web服务",
	RunE: func(c *cobra.Command, args []string) error {
		// 从Command中获取服务容器
		container := c.GetContainer()
		// 从服务容器中获取kernel的服务实例
		//mustmake返回的是个NewInstance,也就是一个实例化函数,是一个函数!!
		//函数是可以直接用的,这里放回的是一个实例化函数
		//kernelServicetemp是个啥?是一个*kernel.HadeKernelService
		//这里比较抽象,什么含义具体?
		//首先调用了container的MustMake方法,这个方法底层会调用make方法,make方法首先会从服务池中寻找是否已经有这个服务的provider
		//如果没有,说明是没有将服务bind入container中,那对不起你用不了,给你报错
		//如果有这个provider,说明容器中有你的注册函数,所以要想使用服务必须先绑定bind
		//然后去instance中找你的key,这里传入的hashkey是kernelkey,如果有kernelkey对应的服务实例,就返回这个实例,
		//没有就调用该key对应的newInstance方法,并且存入instance[],再返回这个服务实例
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		//fmt.Println("kernelServicetemp是个啥?") fmt.Println(reflect.TypeOf(kernelServicetemp))

		// 从kernel服务实例中获取引擎,因为他本来就是个带有engine的结构体
		core := kernelService.HttpEngine()

		// 创建一个Server服务,并将core传入
		server := &http.Server{
			Handler: core,
			Addr:    ":8888",
		}

		// 这个goroutine是启动服务的goroutine
		go func() {
			server.ListenAndServe()
		}()

		// 当前的goroutine等待信号量
		quit := make(chan os.Signal)
		// 监控信号：SIGINT, SIGTERM, SIGQUIT
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		// 这里会阻塞当前goroutine等待信号
		<-quit

		// 调用Server.Shutdown graceful结束
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}

		return nil
	},
}
