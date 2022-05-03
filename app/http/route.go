package http

import (
	"github.com/sunranlike/hade/app/http/module/demo"
	"github.com/sunranlike/hade/framework/gin"
)

// Routes 绑定业务层路由
//调用注册函数
func Routes(r *gin.Engine) {
	//配置文件静态路由
	r.Static("/dist/", "./dist/")
	//注册路由
	demo.Register(r)
}
