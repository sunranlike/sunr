package http

import (

	"github.com/sunranlike/hade/framework/gin"
	"github.com/sunranlike/hade/provider/demo"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {

	r.Static("/dist/", "./dist/")

	demo.Register(r)
}
