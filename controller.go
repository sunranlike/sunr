package main

import (
	"context"
	"fmt"
	"github.com/sunranlike/hade/framework"
	"time"
)

// FooControllerHandler 这个函数就是我们的处理器函数,是我们逻辑的主要实现模块,我们要干嘛就是在这里实行
//在route.go文件中我们吧foo和FooControllerHandler绑定到了一起所以使用"foo"函数就是调用FooControllerHandler 函数
//然后我们又在main中吧注册这个(调用core.get把core和foo->FooControllerHandler绑定到主函数的core中)

func FooControllerHandler(c *framework.Context) error {
	//return ctx.Json(200, map[string]interface{}{
	//	"code": 0,
	//})
	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second)) //1s的ctx,基于ctx
	// 这里记得当所有事情处理结束后调用 cancel，告知 durationCtx 的后续 Context 结束
	defer cancel()

	// 这个 channal 负责通知结束
	finish := make(chan struct{}, 1)
	// 这个 channel 负责通知 panic 异常
	panicChan := make(chan interface{}, 1)

	go func() {
		// 这里增加异常处理,否则出大问题
		defer func() {
			if p := recover(); p != nil {
				//recover是一个内建的函数，可以让进入panic状态的goroutine恢复过来。
				//recover仅在延迟函数中有效。在正常的执行过程中，
				//调用recover会返回nil，并且没有其它任何效果。如
				//果当前的goroutine陷入panic状态，调用recover可以捕获到panic的输入值，并且恢复正常的执行。

				panicChan <- p
			}
		}()
		// 这里做具体的业务
		time.Sleep(10 * time.Second)
		c.IJson("ok").ISetStatus(200)
		//
		// 新的 goroutine 结束的时候通过一个 finish 通道告知父 goroutine
		finish <- struct{}{}
	}()

	select {
	// 监听 panic
	case <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		//TODO
		c.IJson("panic").ISetStatus(500)
		// 监听结束事件
	case <-finish:
		//TODO
		fmt.Println("finish")
		// 监听超时事件
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()

		c.IJson("time out").ISetStatus(500) //已经超时了
		c.SetHasTimeout()                   //告诉大家和这个已经超时了
	}
	return nil
}
