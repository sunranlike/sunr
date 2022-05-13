package command

import (
	"context"
	"fmt"
	"github.com/erikdubbelboer/gspt"
	"github.com/pkg/errors"
	"github.com/sevlyar/go-daemon"
	"github.com/sunranlike/hade/framework"
	"github.com/sunranlike/hade/framework/cobra"
	"github.com/sunranlike/hade/framework/contract"
	"github.com/sunranlike/hade/framework/util"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var appDaemon = false

//　Daemon程序是一直运行的服务端程序，又称为守护进程。
//通常在系统后台运行，没有控制终端不与前台交互，Daemon程序一般作为系统服务使用。
//Daemon是长时间运行的进程，通常在系统启动后就运行，在系统关闭时才结束。
//一般说Daemon程序在后台运行，是因为它没有控制终端，无法和前台的用户交互。
//Daemon程序一般都作为服务程序使用，等待客户端程序与它通信。我们也把运行的Daemon程序称作守护进程。

var appAddress = ""

// AppCommand 是命令行参数第一级为 app 的命令，它没有实际功能，只是打印帮助文档

// initAppCommand 初始化app命令和其子命令
func initAppCommand() *cobra.Command {
	//首先调用Flags()方法,使其成为一个FlagSet,然后在调用spf13的另一个库:pflag
	//pflag的BoolVarP给这个命令添加了一个flag
	appStartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", false, "start app daemon")
	appStartCommand.Flags().StringVar(&appAddress, "address", "", "设置app启动的地址，默认为:8888")

	appCommand.AddCommand(appRestartCommand)
	appCommand.AddCommand(appStateCommand)
	appCommand.AddCommand(appStopCommand)
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

// 启动AppServer, 这个函数会将当前goroutine阻塞
//这个方法主要是为了其他command使用的，主要是给其他命令提供一个优雅关闭的选项，
//而非是在收到客户端ctrl c之后立马杀死服务、

func startAppServe(server *http.Server, c framework.Container) error {
	// 这个goroutine是启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT,这三个信号都是客户端发起的终止操作
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	closeWait := 5
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	if configService.IsExist("app.close_wait") {
		closeWait = configService.GetInt("app.close_wait")
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(closeWait)*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		return err
	}
	return nil
}

// appStartCommand 启动一个Web服务,web服务是框架级服务.
//我们讨论了 start 的关键设计，再回头梳理一遍这个命令的实现步骤：
//一:从四个方式获取参数 appAddress
//二:获取参数 daemon
//三:确认 runtime 目录和 PID 文件存在
//四:确认 log 目录的 log 文件存在
//五:判断是否是 daemon 方式。
//            如果是，就使用 go-daemon 来启动一个子进程；
//            如果不是，直接进行后续调用
//六:使用 gspt 来设置当前进程名称
//七:启动 app 服务
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

		// 从kernel服务实例中获取引擎,因为他本来就是个带有engine的结构体
		core := kernelService.HttpEngine()

		//对appAddress判断:
		//一:直接从命令行参数获取 address 参数，
		//二:是从环境变量 ADDRESS 中获取，
		//三:然后是从配置文件中获取配置项 app.address，
		//最后如果以上三个方式都没有设置，就使用默认值:8888
		if appAddress == "" {
			envService := container.MustMake(contract.EnvKey).(contract.Env)
			if envService.Get("ADDRESS") != "" {
				appAddress = envService.Get("ADDRESS")
			} else {
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configService.IsExist("app.address") {
					appAddress = configService.GetString("app.address")
				} else {
					appAddress = ":8888"
				}
			}
		}

		// 创建一个Server服务
		server := &http.Server{
			Handler: core,
			Addr:    appAddress,
		}

		appService := container.MustMake(contract.AppKey).(contract.App)
		//配置pid和log   "app.pid"  "app.log"
		pidFolder := appService.RuntimeFolder()
		if !util.Exists(pidFolder) {
			if err := os.MkdirAll(pidFolder, os.ModePerm); err != nil {
				return err
			}
		}
		serverPidFile := filepath.Join(pidFolder, "app.pid")
		logFolder := appService.LogFolder()
		if !util.Exists(logFolder) {
			if err := os.MkdirAll(logFolder, os.ModePerm); err != nil {
				return err
			}
		}
		// 应用日志
		serverLogFile := filepath.Join(logFolder, "app.log")
		currentFolder := util.GetExecDirectory()
		// daemon 模式,appDeamon是个flag,在调用app start的时候就会声明是否守护进程,默认false
		if appDaemon {
			// 创建一个Context.deamon守护进程的核心就在于这个ctx
			cntxt := &daemon.Context{
				// 设置pid文件
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				// 设置日志文件
				LogFileName: serverLogFile,
				LogFilePerm: 0640,
				// 设置工作路径
				WorkDir: currentFolder,
				// 设置所有设置文件的mask，默认为750
				Umask: 027,
				// 子进程的参数，按照这个参数设置，子进程的命令为 ./hade app start --daemon=true
				Args: []string{"", "app", "start", "--daemon=true"},
			}
			// 启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
			//Reborn 理解成 fork，当调用这个函数的时候，父进程会继续往下走，但是返回值 d 不为空，它的信息是子进程的进程号等信息。
			//而子进程会重新运行对应的命令，再次进入到 Reborn 函数的时候，返回的 d 就为 nil。所以在 Reborn 的后面，
			//我们让父进程直接 return，而让子进程继续往后进行操作，这样就达到了 fork 一个子进程的效果了。
			//既:reborn结束父进程
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				// 父进程直接打印启动成功信息，不做任何操作
				fmt.Println("app启动成功，pid:", d.Pid)
				fmt.Println("日志文件:", serverLogFile)
				return nil
			}
			defer cntxt.Release()
			// 子进程执行真正的app启动操作
			fmt.Println("deamon started")
			//对于启动的进程，我们一般都希望能自定义它的进程名称。
			//gspt可以直接修改进程的名称
			gspt.SetProcTitle("hade app")
			if err := startAppServe(server, container); err != nil {
				fmt.Println(err)
			}
			return nil
		}

		// 非deamon模式，直接执行
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := ioutil.WriteFile(serverPidFile, []byte(content), 0644)
		if err != nil {
			return err
		}
		gspt.SetProcTitle("hade app")

		fmt.Println("app serve url:", appAddress)
		if err := startAppServe(server, container); err != nil {
			fmt.Println(err)
		}
		return nil
	},
}

// 重新启动一个app服务
//同其他命令一样，这里再梳理一下判断旧进程存在之后详细的实现步骤，如果存在：
//    发送 SIGTERM 信号
//    循环 2*closeWait 次数，每秒执行一次查询进程是否已经结束
//    如果某次查询进程已经结束，或者等待 2*closeWait 循环结束之后，再次查询一次进程
//    如果还未结束，返回进程结束失败
//    如果已经结束，将 PID 文件清空，启动新进程
var appRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重新启动一个app服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 文件中获取GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		if !util.Exists(serverPidFile) { //如果不存在直接以守护进程启动这个app,为什么要守护进程
			appDaemon = true
			// 直接daemon方式启动apps
			return appStartCommand.RunE(c, args)
		}

		content, err := ioutil.ReadFile(serverPidFile) //读取文件
		if err != nil {
			return err
		}

		if content != nil && len(content) != 0 {
			pid, err := strconv.Atoi(string(content)) //strconv强制转换为int
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) { //这个util检查是否存在对应的进程,若存在则杀死,因为我们要重启
				// 杀死进程,
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					return err
				}

				// 获取closeWait,轮询closeWait秒 再次检查是否存在这个进程,确保这个进程真正被杀死
				//因为我们这些服务公用的是一个端口
				closeWait := 5
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configService.IsExist("app.close_wait") {
					closeWait = configService.GetInt("app.close_wait")
				}

				// 确认进程已经关闭,每秒检测一次， 最多检测closeWait * 2秒
				for i := 0; i < closeWait*2; i++ {
					if util.CheckProcessExist(pid) == false {
						break
					}
					time.Sleep(1 * time.Second)
				}

				// 如果进程等待了2*closeWait之后还没结束，返回错误，不进行后续的操作
				if util.CheckProcessExist(pid) == true {
					fmt.Println("结束进程失败:"+strconv.Itoa(pid), "请查看原因")
					return errors.New("结束进程失败")
				}
				if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
					return err
				}

				fmt.Println("结束进程成功:" + strconv.Itoa(pid))
			}
		}

		appDaemon = true
		// 直接daemon方式启动apps
		return appStartCommand.RunE(c, args)
	},
}

// 停止一个已经启动的app服务
//同样实现步骤也很清晰，获取 PID 文件内容之后，判断如果有 PID 文件且有内容再继续，否则什么都不做，之后就是：
//一:将内容转换为 PID 的 int 类型，转换失败则什么都不做
//二:直接给这个 PID 进程发送 SIGTERM 信号
//三:将 PID 文件内容清空
var appStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止一个已经启动的app服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 从文件中get pid文件.
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) != 0 { //没有app.id文件就说明没有这个服务,直接返回
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			// 有这个文件,对这个文件中的pid,发送SIGTERM命令,杀死程序
			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}
			//644是linux的文件权限,644说明自己可以读取和写入,但是不能执行(他也不是可执行文件)
			//4的含义是只可读取,
			//三位数字第一个是自己,第二个是同用户组,第三个是所有人,既我自己是6(rw),同用户组人是4(r),其他人也是4(r)
			//同理还有其他权限比如说777,7就是二进制的111,对应的就是rwx的111,既可读可写可执行,那么777就是很宽泛的权限,所有的人都可以读写执行
			if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil { //写入文件记录
				return err
			}
			fmt.Println("停止进程:", pid)
		}
		return nil
	},
}

//获取启动的app的pid
//获取 PID 文件内容之后，做判断，如果有 PID 文件且有内容就继续操作，否则无文件说明无进程,返回无进程；
//如果检测到本地有app id 文件：
//    将内容转换为 PID 的 int 类型，转换失败视为无进程；
//    使用 signal 0 确认这个进程是否存在，存在返回结果有进程，不存在返回结构无进程。
var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "获取启动的app的pid",
	RunE: func(c *cobra.Command, args []string) error {
		//获取容器
		container := c.GetContainer()
		//创建爱你app目录服务.
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 在runtimeFolder下获取app.id文件并判断文件是否存在
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		//如果文件存在说明
		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				fmt.Println("app服务已经启动, pid:", pid)
				return nil
			}
		}
		fmt.Println("没有app服务存在")
		return nil
	},
}
