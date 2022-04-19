package framework

import (
	"log"
	"net/http"
)

// Core 框架核心结构,core这个结构体是干嘛的?其实就是一个处理器(handler)表,有一个url对应的func
type Core struct {
	router map[string]ControllerHandler //暂时只有一个map,
}

// NewCore 初始化框架核心结构,Newcore干嘛了?其实就是是:返回一个core结构体,这个core结构体
//有一个router的map.key是string,value是一个ControllerHandler,是个函数,
//这也是一种设计模式吧,key是string,对用的value是一个函数,这叫做函数工厂?
func NewCore() *Core {
	return &Core{router: map[string]ControllerHandler{}} //字符串映射到函数方程
}

// Get Get方法干嘛了?就是把参数作为key;value存入map中,相当于注册,为什么叫这个?
func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[url] = handler
}

// 实现core的Handler接口,既实现ServeHTTP方法，如果你不自己实现就会调用默认的serveHttp方法
//该函数不需要主动调用,只需要调用ListenAndServe方法,后台就会自动调用这个servehttp方法
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core.serveHTTP")
	log.Println(request.URL, request.Method, request.Body)
	ctx := NewContext(request, response)

	// 一个简单的路由选择器，这里直接写死为测试路由foo
	router := c.router["foo"] //rotuer 是什么?是core的一个map,map映射的是方程,所以这里router最后是个函数
	if router == nil {
		log.Println("no foo func")
		return
	}

	log.Println("core.router")

	router(ctx) //执行这个取出来的函数,当然要把ctx传入
}

//在源码中Handler 接口实际上只有一个函数 就是ServerHTTP方法，所以我们自己写一个ServeHttp就代表我们
//使用自己的Handler，（实际上是一个ServerHttp方法）
//type Handler interface {
//	ServeHTTP(ResponseWriter, *Request)
//}
