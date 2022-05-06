package console

import (
	"github.com/sunranlike/hade/framework"
	"github.com/sunranlike/hade/framework/cobra"
	"github.com/sunranlike/hade/framework/command"
)

// RunCommand  初始化根Command并运行
//该函数绑定了一个container,我们的命令行都会围绕container

func RunCommand(container framework.Container) error {
	// 根Command
	//cobra.Command是一个结构体,这个结构体在赋值的时候:会执行RunE的字段
	//这是一个根Command,我们的业务,框架服务想要添加新的command,都会是他的子command结构
	var rootCmd = &cobra.Command{
		// 定义根命令的关键字
		Use: "hade",
		// 简短介绍
		Short: "hade 命令",
		// 根命令的详细介绍
		Long: "hade 框架提供的命令行工具，使用这个命令行工具能很方便执行框架自带命令，也能很方便编写业务命令",
		// 根命令的执行函数,赋值的时候就会执行这个方法
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag() //初始化一下
			return cmd.Help()         //打印一个help
		},
		// 不需要出现cobra默认的completion子命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	// 将我们的参数container作为根Command的服务容器,实现了服务传入Command的功能
	//这样rootCmd就有一个container结构体了,这个结构体可以被他的子command使用get命令获取得到
	rootCmd.SetContainer(container)

	// 执行AddKernelCommands 绑定框架的命令
	command.AddKernelCommands(rootCmd)

	// 绑定业务的命令
	AddAppCommand(rootCmd)

	// 执行RootCommand
	return rootCmd.Execute()
}

// 绑定业务的命令
func AddAppCommand(rootCmd *cobra.Command) {
	//  demo 例子
	//rootCmd.AddCommand(demo.InitFoo())
	//rootCmd.AddDistributedCronCommand("foo_func_for_test", "* * * * * *", demo.FooCommand, 2*time.Second)
}
