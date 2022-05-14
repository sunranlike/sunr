package command

import (
	"fmt"
	"github.com/sunranlike/sunr/framework/cobra"
	"github.com/sunranlike/sunr/framework/contract"
)

// helpCommand show current envionment
//一个demo的框架级命令
var FrameworkDemoCommand = &cobra.Command{
	Use:   "demo",
	Short: "demo for framework",
	Run: func(c *cobra.Command, args []string) {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		fmt.Println("app base folder:", appService.BaseFolder())
	},
}
