package cobra

import (
	"github.com/robfig/cron/v3"
	"github.com/sunranlike/sunr/framework"
	"log"
)

//SetContainer 设置服务容器,将外部的服务容器container传入Command结构中,实现了容器的传递.
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

//GetContainer 获取自己的容器,并且因为只有根节点Command有这个节点
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}

// CronSpec 保存Cron命令的信息，用于展示
type CronSpec struct {
	Type        string
	Cmd         *Command
	Spec        string
	ServiceName string
}

func (c *Command) SetParantNull() {
	c.parent = nil
}

// AddCronCommand 是用来创建一个Cron任务的
func (c *Command) AddCronCommand(spec string, cmd *Command) {
	// cron结构是挂载在根Command上的
	root := c.Root()
	if root.Cron == nil { //先判断有没有初始化cron
		// 初始化cron
		root.Cron = cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
		root.CronSpecs = []CronSpec{}
	}
	// 增加说明信息
	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type: "normal-cron",
		Cmd:  cmd,
		Spec: spec,
	})

	// 制作一个rootCommand
	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd //这里创建了一个新的cronCmd指向了父cmd
	cronCmd.args = []string{}
	cronCmd.SetParantNull()
	cronCmd.SetContainer(root.GetContainer())

	// 增加调用函数
	root.Cron.AddFunc(spec, func() {
		// 如果后续的command出现panic，这里要捕获,这个addFunc应该是在所有的函数执行完之后
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			// 打印出err信息
			log.Println(err)
		}
	})
}
