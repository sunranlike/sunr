package command

import "github.com/sunranlike/hade/framework/cobra"

// AddKernelCommands will add all command/* to root command
func AddKernelCommands(root *cobra.Command) {
	//root.AddCommand(FrameworkDemoCommand)
	root.AddCommand(initEnvCommand())

	root.AddCommand(initCronCommand())

	root.AddCommand(initBuildCommand())
	root.AddCommand(initConfigCommand())
	//
	//// app
	//传入框架级别服务:目录app服务,到rootcmd中
	root.AddCommand(initAppCommand())

}
