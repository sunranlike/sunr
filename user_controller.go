package main

import (
	_ "github.com/sunranlike/hade/framework"
	"github.com/sunranlike/hade/framework/gin"
	"time"
)

func UserLoginController(c *gin.Context) {
	//c.Json("ok, UserLoginController").SetStatus(200)
	foo, _ := c.DefaultQueryString("foo", "def")
	// 等待10s才结束执行
	time.Sleep(10 * time.Second)
	// 输出结果
	c.ISetOkStatus().IJson("ok, UserLoginController: " + foo)

}
