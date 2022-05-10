package http

import (
	"github.com/sunranlike/hade/app/http/module/demo"
	"github.com/sunranlike/hade/framework/contract"
	"github.com/sunranlike/hade/framework/gin"
	ginSwagger "github.com/sunranlike/hade/framework/middleware/gin-swagger"
	"github.com/sunranlike/hade/framework/middleware/gin-swagger/swaggerFiles"
	"github.com/sunranlike/hade/framework/middleware/static"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {
	container := r.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	// /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))
	//r.Use(cors.Default())

	// 如果配置了swagger，则显示swagger的中间件
	if configService.GetBool("app.swagger") == true {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	// 动态路由定义
	demo.Register(r)
}

// 注册路由匹配规则,什么样的uri匹配什么样的Controller
//路由匹配的需求:
//1:静态路由匹配:一个url对应一个handler
//2:静态路由分组匹配:其实就是1的进阶版,可以嵌套
//3:动态路有字段匹配

//func registerRouter(core *gin.Engine) { //
//	// 静态路由+HTTP方法匹配
//	//get 实际上是注册url： UserLoginController 到map上 在serve的时候会根据与这些个get规则执行对应的方法与Handler
//	//根据http  head 的method 执行对应的Handler
//	//在核心业务逻辑 UserLoginController 之外，封装一层 TimeoutHandler
//	//get方法将url网址和handler处理器 绑定一起,但其实底层都是一个map对应一个字典树
//	//core.GET("/user/login", middleware.Cost(), UserLoginController)
//
//	//上一步core绑定了"/user/login"这个urlUI对应的handler,这个就是所谓的静态路由,一个url对应一个handler
//	//效率不高,但是我们可以通过路由组group功能实现批量匹配,提高效率
//	//我们通过Group批量匹配功能,其实也就是路由组，实现批量通用前缀
//
//	//整体流程就是1通过core的Group方法创建router group路由组，然后在使用这个路由组的GEt、PUT、
//
//	//subjectApi1 := core.Group("/subject") //调用Core的Group方法(有一个结构体也叫group),
//	//只要匹配到url匹配到/subject,就会进入这个group
//
//	//这里为什么要单独弄一个花括号?
//	//语法上，我们习惯（也是官方建议）将一组路由放在一个代码块中，在结构上保持独立。
//	//但这个代码块不是必要的。!!!!!
//
//	//{
//	//	// 动态路由"/:id"
//	//	//只要匹配到"/:id" 就会执行delete put get 三个方法
//	//	subjectApi1.DELETE("/:id", SubjectDelController) //Group调用get/put 方法本质上还是调用的core的,只不过修饰了一个中间件
//	//	subjectApi1.PUT("/:id", SubjectUpdateController)
//	//	subjectApi1.GET("/:id", SubjectGetController)
//	//	subjectListApi := subjectApi1.Group("/list") //在Api1上再 在添加一个Api，也就是一个新的group，路由进入新的api也要符合上一层的路由规则
//	//	//分组的关键，这样子才能嵌套，必须要使用上层的group作为右值
//	//	//这样的左值在使用get方法才会形成嵌套
//	//	{
//	//		subjectListApi.GET("/all", SubjectListController)
//	//	}
//	//
//	//	subjectInnerApi := subjectApi1.Group("/info")
//	//	{
//	//		subjectInnerApi.GET("/name", SubjectNameController)
//	//	}
//	//}
//
//	//core.PrintRouter()
//	//经过这个group分组,最终是一个
//
//}
