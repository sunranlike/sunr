package command

import "github.com/sunranlike/hade/framework/cobra"

// AddKernelCommands will add all command/* to root command
func AddKernelCommands(root *cobra.Command) {
	//root.AddCommand(FrameworkDemoCommand)
	root.AddCommand(initEnvCommand())
	//root.AddCommand(envCommand)
	//root.AddCommand(deployCommand)
	//
	//// cron
	root.AddCommand(initCronCommand())
	//// cmd
	//cmdCommand.AddCommand(cmdListCommand)
	//cmdCommand.AddCommand(cmdCreateCommand)
	//root.AddCommand(cmdCommand)
	//
	//// build
	//buildCommand.AddCommand(buildSelfCommand)
	//buildCommand.AddCommand(buildBackendCommand)
	//buildCommand.AddCommand(buildFrontendCommand)
	//buildCommand.AddCommand(buildAllCommand)
	//root.AddCommand(buildCommand)
	//
	//// app
	//传入框架级别服务:目录app服务,到rootcmd中
	root.AddCommand(initAppCommand())
	//
	//// dev
	//root.AddCommand(initDevCommand())
	//
	//// middleware
	//middlewareCommand.AddCommand(middlewareAllCommand)
	//middlewareCommand.AddCommand(middlewareAddCommand)
	//middlewareCommand.AddCommand(middlewareRemoveCommand)
	//root.AddCommand(middlewareCommand)
	//
	//// swagger
	//swagger.IndexCommand.AddCommand(swagger.InitServeCommand())
	//swagger.IndexCommand.AddCommand(swagger.GenCommand)
	//root.AddCommand(swagger.IndexCommand)
	//
	//// provider
	//providerCommand.AddCommand(providerListCommand)
	//providerCommand.AddCommand(providerCreateCommand)
	//root.AddCommand(providerCommand)
	//
	//// new
	//root.AddCommand(initNewCommand())
}
