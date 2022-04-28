package main

import (
	"github.com/sunranlike/hade/framework/gin"
	"github.com/sunranlike/hade/provider/demo"
)

func SubjectAddController(c *gin.Context) {
	c.IJson("ok, SubjectAddController").ISetStatus(200)

}

// 对应路由 /subject/list/all
func SubjectListController(c *gin.Context) {
	// 获取 demo 服务实例
	demoService := c.MustMake(demo.Key).(demo.Service)

	// 调用服务实例的方法
	foo := demoService.GetFoo()

	// 输出结果
	c.ISetOkStatus().IJson(foo)
}

func SubjectDelController(c *gin.Context) {
	c.IJson("ok, SubjectDelController").ISetStatus(200)

}

func SubjectUpdateController(c *gin.Context) {
	c.IJson("ok, SubjectUpdateController").ISetStatus(200)

}

func SubjectGetController(c *gin.Context) {
	c.IJson("ok, SubjectGetController").ISetStatus(200)

}

func SubjectNameController(c *gin.Context) {
	c.IJson("ok, SubjectNameController").ISetStatus(200)

}
