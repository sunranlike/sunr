package middleware

import (
	"github.com/sunranlike/hade/framework/gin"
)

// Recovery :recovery机制，将协程中的函数异常进行捕获
//经过recover机制的handler会更加健壮,具有处理下层panci的能力
func Recovery() gin.HandlerFunc {
	// 使用函数回调
	return func(c *gin.Context) {
		// 核心在增加这个recover机制，捕获c.Next()出现的panic
		defer func() {
			if err := recover(); err != nil {
				c.IJson("err").ISetStatus(500)
			}
		}()
		//time.Sleep(3*time.Second)
		// 使用next执行具体的业务逻辑
		c.Next()

	}
}
