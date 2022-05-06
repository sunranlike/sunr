package command

import (
	"fmt"
	"github.com/sunranlike/hade/framework/cobra"
	"log"
	"os/exec"
)

func initBuildCommand() *cobra.Command { //返回指针
	buildCommand.AddCommand(buildSelfCommand)
	buildCommand.AddCommand(buildBackendCommand)
	buildCommand.AddCommand(buildFrontendCommand)
	buildCommand.AddCommand(buildAllCommand)
	return buildCommand
}

var buildCommand = &cobra.Command{ //本来就是个指针
	Use:   "build",
	Short: "编译相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 { //如果没有长度,参数长度为0 就是 hade build 没有下一步了,就会提示help
			c.Help()
		}
		return nil
	},
}
var buildSelfCommand = &cobra.Command{ //应该是为了自编译>
	Use:   "self",
	Short: "编译hade命令",
	RunE: func(c *cobra.Command, args []string) error { //判断有无安装go
		path, err := exec.LookPath("go")
		if err != nil {
			log.Fatalln("hade go: please install go in path first")
		}

		cmd := exec.Command(path, "build", "-o", "hade", "./") //编译
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("go build error:")
			fmt.Println(string(out))
			fmt.Println("--------------")
			return err
		}
		fmt.Println("build success please run ./hade direct")
		return nil
	},
}
var buildBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "使用go编译后端",
	RunE: func(c *cobra.Command, args []string) error {
		return buildSelfCommand.RunE(c, args) //编译后端.说就是直接调用自编译
	},
}

// 打印前端的命令
var buildFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "使用npm编译前端",
	RunE: func(c *cobra.Command, args []string) error {
		// 获取path路径下的npm命令,你的path环境下一定要有
		path, err := exec.LookPath("npm") //在path下面是有的
		if err != nil {
			log.Fatalln("请安装npm在你的PATH路径下")
		}

		// 执行npm run build
		cmd := exec.Command(path, "run", "build")
		// 将输出保存在out中
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("=============  前端编译失败 ============")
			fmt.Println(string(out))
			fmt.Println("=============  前端编译失败 ============")
			return err
		}
		// 打印输出
		fmt.Print(string(out))
		fmt.Println("=============  前端编译成功 ============")
		return nil
	},
}

var buildAllCommand = &cobra.Command{
	Use:   "all",
	Short: "同时编译前端和后端",
	RunE: func(c *cobra.Command, args []string) error {
		err := buildFrontendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		err = buildBackendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		return nil
	},
}
