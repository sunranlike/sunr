// Copyright 2021 jianfengye.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"github.com/sunranlike/hade/framework"
)

//func (ctx *Context) BaseContext() context.Context {
//	return ctx.Request.Context()
//}

// Bind engine实现container的绑定封装
//？这是个什么逻辑？我们在实现它的时候在调用它？
//其实engine没有实现Bind方法，他只是调用了container字段的Bind方法，
//其实也是合法的，因为规定container字段的value是一个实现了framework.Container的字段，
//那么container字段必定有一个Bind方法
//实际上我们在调用gin.New的时候，给这个container字段返回的是hadecontainer，这个歌结构体是实现了
//framerwork.Container接口的，也就是他又bind和make类方法，自然就可以用
func (engine *Engine) Bind(provider framework.ServiceProvider) error {
	return engine.container.Bind(provider)
}

// IsBind 关键字凭证是否已经绑定服务提供者
func (engine *Engine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}

// context 实现container的几个封装

// 实现make的封装
func (ctx *Context) Make(key string) (interface{}, error) {
	return ctx.container.Make(key)
}

// 实现mustMake的封装
func (ctx *Context) MustMake(key string) interface{} {
	return ctx.container.MustMake(key)
}

// 实现makenew的封装
func (ctx *Context) MakeNew(key string, params []interface{}) (interface{}, error) {
	return ctx.container.MakeNew(key, params)
}
