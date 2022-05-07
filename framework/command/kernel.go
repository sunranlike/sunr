package command

import "github.com/sunranlike/hade/framework/cobra"

// AddKernelCommands will add all command/* to root command
func AddKernelCommands(root *cobra.Command) {

	root.AddCommand(initAppCommand())

	root.AddCommand(initEnvCommand())

	root.AddCommand(initCronCommand())

	root.AddCommand(initDevCommand())

	root.AddCommand(initBuildCommand())

}
