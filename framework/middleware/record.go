package middleware

import (
	"coredemo/framework"
	"log"
	"time"
)

func RecordRequsstTime() framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) error {
		// 获取开始时间
		startT := time.Now()
		// 输出请求URI

		// 执行其他中间件和函数处理
		c.Next()
		// 获取处理时长
		time.Sleep(1 * time.Second)
		c.Json("休息1s哈哈哈")
		tc := time.Since(startT)
		log.Println(tc)
		return nil
	}
}
