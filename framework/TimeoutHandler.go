package framework

import (
	"context"
	"fmt"
	"log"
	"time"
)

// TimeoutHandler 这是一个装饰器函数，参数和返回值都是一个ControllerHandler结构体
func TimeoutHandler(fun ControllerHandler, d time.Duration) ControllerHandler {
	// 使用函数回调
	return func(c *Context) error { //这就是ControllerHandler的函数签名

		//两个chan，一个finish一个panic chan
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		// 执行业务逻辑前预操作：初始化超时 context
		//Without其实也是一个装饰器模式
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		//调用request的WithContext，把我们的durationCtx绑定到c上，
		//c是我们的传入的参数ctx，我们要对着ctx修饰一个duration，这样调用我们的调用者
		//就会得到一个带有过期时间（当然也需要作为参数传入）的request结构体
		c.request.WithContext(durationCtx)

		//干嘛的，执行具体逻辑：也就是执行传入的handler，
		go func() {
			defer func() { //异常捕获
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 执行具体的业务逻辑
			fun(c)

			finish <- struct{}{}
		}()
		// 执行业务逻辑后操作：即
		//1监听业务逻辑代码是否报错panicChan
		//2监听业务逻辑是否执行完毕finish
		//3持续时间是否完成
		select {
		case p := <-panicChan:
			log.Println(p)
			c.responseWriter.WriteHeader(500)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.responseWriter.Write([]byte("time out"))
		}
		return nil
	}
}
