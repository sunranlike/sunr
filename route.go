package main

//route是干嘛的?其实就是路由功能,将一些方法注册到core的,把参数作为key;value存入map中
import (
	"coredemo/framework"
	"time"
)

//吧foo->FooControllerHandler 这个映射与core绑定到一起
//func registerRouter(core *framework.Core) {
//	// 设置控制器,将foo和FooControllerHandler这个我们自己写的函数绑定起来
//	core.Get("foo", FooControllerHandler)
//}

// 注册路由匹配规则,什么样的uri匹配什么样的Controller
func registerRouter(core *framework.Core) {
	// 静态路由+HTTP方法匹配
	//get 实际上是注册url： UserLoginController 到map上 在serve的时候能够
	//根据http  head 的method 执行对应的Handler
	//在核心业务逻辑 UserLoginController 之外，封装一层 TimeoutHandler
	core.Get("/user/login",
		framework.TimeoutHandler(UserLoginController, time.Second))

	// 批量通用前缀,这也是Group的一个批量匹配规则
	subjectApi := core.Group("/subject")
	//只要匹配到url匹配到/subject,就会进入这个group
	{
		// 动态路由
		//只要匹配到"/:id" 就会全部执行
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Put("/:id", SubjectUpdateController)

		subjectApi.Get("/:id", SubjectGetController)
		subjectApi.Get("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.Get("/name", SubjectNameController)
		}
	}
}
