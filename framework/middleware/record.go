package middleware

import (
	_ "github.com/sunranlike/sunr/framework"
	"github.com/sunranlike/sunr/framework/gin"
	"log"
	"time"
)

func Cost() gin.HandlerFunc {
	// 使用函数回调
	return func(c *gin.Context) {
		// 获取开始时间
		startT := time.Now()
		// 输出请求URI
		log.Printf("api uri start: %v", c.Request.RequestURI)
		// 执行其他中间件和函数处理
		c.Next()
		// 获取处理时长
		//time.Sleep(1 * time.Second)
		//c.Json("休息1s哈哈哈")
		tc := time.Since(startT)
		log.Printf("api uri end: %v, cost: %v", c.Request.RequestURI, tc)

	}
}
