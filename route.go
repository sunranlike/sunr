package main

//route是干嘛的?其实就是路由功能,将一些方法注册到core的,把参数作为key;value存入map中
import "coredemo/framework"

//吧foo->FooControllerHandler 这个映射与core绑定到一起
func registerRouter(core *framework.Core) {
	// 设置控制器,将foo和FooControllerHandler这个我们自己写的函数绑定起来
	core.Get("foo", FooControllerHandler)
}
