package framework

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Core 框架核心结构,core这个结构体是干嘛的?他不仅包含了一个map,并且通过一个Tree结构实现了uri匹配
//

type Core struct {
	router      map[string]*Tree    // 本来是双层的数据结构,后来为了实现匹配查找用了树
	middlewares []ControllerHandler // 从core这边设置的中间件
}

func (core *Core) PrintRouter() {
	fmt.Println(core.router)
	fmt.Println(core.middlewares)
}

// NewCore 初始化框架核心结构,Newcore干嘛了?其实就是是:返回一个core结构体,这个core结构体
//有一个router的map.key是string,value是一个ControllerHandler,是个函数,
//这也是一种设计模式吧,key是string,对用的value是一个函数,这叫做函数工厂?应该不是，就是一个映射关系
func NewCore() *Core {
	//Core的router原来是一个二级map,经过改造不是map了
	//现在的map是一个
	//实现了tree保存函数：
	tempRouter := map[string]*Tree{} //先声明,然后这个temprouter会作为Core结构体的router这个field的值
	tempRouter["GET"] = NewTree()    //存入这么二级map中
	tempRouter["POST"] = NewTree()
	tempRouter["PUT"] = NewTree()
	tempRouter["DELETE"] = NewTree()
	//以上操作将方法存入
	return &Core{router: tempRouter} //返回这个结构,当然我们只对router做了赋值,而中间件还没有
}

//Get 方法首先将传入的处理器handlers存入c.middlewares这个slice,也就是core的中间件当中
//然后调用调用c.router["GET"] 这个tree结构的AddRouter,这个方法将url网址和处理器 allHandlers 绑定
//这样就可以使得处理器绑定到对应的url,使得
func (c *Core) Get(url string, handlers ...ControllerHandler) {
	//upperUrl := strings.ToUpper(url)
	//c.router["GET"][upperUrl] = handler//注册为入二级map
	//通过变长参数将数据添加入middlewares[] 中间件slice中,未来会被链式调.
	allHandlers := append(c.middlewares, handlers...)

	//将中间件middlerwares添加入路由,绑定router这个map的get方法,同时将url绑定入入这个中间件middlerware.
	//完成了url,get方法,handler三者绑定起来
	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = PUT
func (c *Core) Put(url string, handlers ...ControllerHandler) {
	//upperUrl := strings.ToUpper(url)
	//c.router["PUT"][upperUrl] = handler
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["POST"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = Post
func (c *Core) Post(url string, handlers ...ControllerHandler) {
	//upperUrl := strings.ToUpper(url)
	//c.router["POST"][upperUrl] = handler
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["PUT"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = DELETE
func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	//upperUrl := strings.ToUpper(url)
	//c.router["DELETE"][upperUrl] = handler
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["DELETE"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Group(prefix string) IGroup { //为什么这样写?我在调用方直接用 NewGroup不行吗?
	return NewGroup(c, prefix)
}

// 通过request来匹配路由找到对应的方法，如果没有匹配到，返回nil
func (c *Core) FindRouteNodeByRequest(request *http.Request) *node {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)
	//upperUri := strings.ToUpper(uri)

	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.root.matchNode(uri)
	}
	return nil
}

// 实现core的Handler接口,既实现ServeHTTP方法，如果你不自己实现就会调用默认的serveHttp方法
//该函数不需要主动调用,只需要调用ListenAndServe方法,后台就会自动调用这个servehttp方法
//这个ServeHttp是实际的业务逻辑,
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core.serveHTTP")
	log.Println(request.URL, request.Method)
	ctx := NewContext(request, response) //这个已经不是简单地ctx了而是一个功能丰富的机构提
	//只是名字还和ctx一样
	// 寻找路由
	node := c.FindRouteNodeByRequest(request) //升级为双向的
	if node == nil {
		// 如果没有找到，这里打印日志
		ctx.IJson("not found")
		//ctx.SetStatus(404).Json("not found")
		return
	}

	//rotuer 是什么?是core的一个map,map映射的是一个tree，这个tree是对应的方法
	//寻找路由

	log.Println("core.router")
	// 设置context中的handlers字段
	ctx.SetHandlers(node.handlers)
	params := node.parseParamsFromEndNode(request.URL.Path)
	ctx.SetParams(params)
	// 调用路由函数，如果返回err 代表存在内部错误，返回500状态码
	if err := ctx.Next(); err != nil { //第一次的时候需要调用Next
		ctx.IJson("inner error").ISetStatus(500)
		return
	}
}

// Use 在源码中Handler 接口实际上只有一个函数 就是ServerHTTP方法，所以我们自己写一个ServeHttp就代表我们
//使用自己的Handler，（实际上是一个ServerHttp方法）
// Use 注册中间件,将变长参数middlewares绑定到core的middlewares中
//使用变长参数的好处是可以加好多的参数比如说:
//core.Use(middleware.Recovery(),middleware.RecordRequsstTime())
func (c *Core) Use(middlewares ...ControllerHandler) { //变长参数是,三个...
	c.middlewares = append(c.middlewares, middlewares...) //同样的append也是...
}
