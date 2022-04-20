package framework

import (
	"log"
	"net/http"
	"strings"
)

// Core 框架核心结构,core这个结构体是干嘛的?其实就是一个处理器(handler)表,有一个url对应的func

type Core struct {
	router map[string]*Tree // 二层hashmap
}

// NewCore 初始化框架核心结构,Newcore干嘛了?其实就是是:返回一个core结构体,这个core结构体
//有一个router的map.key是string,value是一个ControllerHandler,是个函数,
//这也是一种设计模式吧,key是string,对用的value是一个函数,这叫做函数工厂?应该不是，就是一个映射关系
func NewCore() *Core {
	//定义二级map的内容
	//getRouter := map[string]ControllerHandler{}
	//postRouter := map[string]ControllerHandler{}
	//putRouter := map[string]ControllerHandler{}
	//deleteRouter := map[string]ControllerHandler{}
	//
	//// 将二级map写入一级map，这样就完成了一个映射关系，外层的四个关键字Get，Post，Put，Delet
	////又各自一一对应一个map，形成了二级map
	//router := map[string]map[string]ControllerHandler{}
	//router["GET"] = getRouter
	//router["POST"] = postRouter
	//router["PUT"] = putRouter
	//router["DELETE"] = deleteRouter
	//
	//return &Core{router: router}
	//实现了tree保存函数：
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{router: router}
}

// Get Get方法干嘛了?就是把参数作为key;value存入二级map中,相当于路由注册：调用这个函数就会生成
//将url和对应的handler注册入map中 映射map就是：
//router["GET"][upperUrl] = handler

func (c *Core) Get(url string, handler ControllerHandler) {
	//upperUrl := strings.ToUpper(url)
	//c.router["GET"][upperUrl] = handler//注册为入二级map
	if err := c.router["GET"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = PUT
func (c *Core) Put(url string, handler ControllerHandler) {
	//upperUrl := strings.ToUpper(url)
	//c.router["PUT"][upperUrl] = handler
	if err := c.router["POST"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = Post
func (c *Core) Post(url string, handler ControllerHandler) {
	//upperUrl := strings.ToUpper(url)
	//c.router["POST"][upperUrl] = handler
	if err := c.router["PUT"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = DELETE
func (c *Core) Delete(url string, handler ControllerHandler) {
	//upperUrl := strings.ToUpper(url)
	//c.router["DELETE"][upperUrl] = handler
	if err := c.router["DELETE"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// 通过request来匹配路由找到对应的方法，如果没有匹配到，返回nil
func (c *Core) FindRouteByRequest(request *http.Request) ControllerHandler {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)
	//upperUri := strings.ToUpper(uri)

	// 查找第一层map，common ok格式的查找
	//先查找第一层，匹配Method，也就是get这些
	//if methodHandlers, ok := c.router[upperMethod]; ok {
	//	// 查找第二层map，也就是取到了对应的method去匹配对应的uri
	//	if handler, ok := methodHandlers[upperUri]; ok {
	//		return handler//返回处理器
	//	}
	//}
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.FindHandler(uri) //估计已经实现了大小写匹配了
	}
	return nil
}

// 实现core的Handler接口,既实现ServeHTTP方法，如果你不自己实现就会调用默认的serveHttp方法
//该函数不需要主动调用,只需要调用ListenAndServe方法,后台就会自动调用这个servehttp方法
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core.serveHTTP")
	log.Println(request.URL, request.Method, request.Body)
	ctx := NewContext(request, response) //这个已经不是简单地ctx了而是一个功能丰富的机构提
	//只是名字还和ctx一样

	// 一个简单的路由选择器，这里直接写死为测试路由foo
	//rotuer 是什么?是core的一个map,map映射的是一个函数,所以这里router最后是个函数
	router := c.FindRouteByRequest(request)
	if router == nil {
		ctx.Json(404, "not found")
		return
	}

	log.Println("core.router")

	if err := router(ctx); err != nil {
		ctx.Json(500, "inner error")
		return
	}
	//执行这个取出来的函数,当然要把ctx传入
}

//在源码中Handler 接口实际上只有一个函数 就是ServerHTTP方法，所以我们自己写一个ServeHttp就代表我们
//使用自己的Handler，（实际上是一个ServerHttp方法）
//type Handler interface {
//	ServeHTTP(ResponseWriter, *Request)
//}
