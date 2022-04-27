package framework

/*
这个是自己写的context,为的就是方便使用,当然自己的ctx需要实现官方包的ctx的方法
然而这个方法并不需要你手动实现,你只需要直接用*http.Request包的context
怎么用:使用BaseContext()函数,他直接
*/
/*


import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// 自定义的 Context,使用request的ctx作为基本的ctx
//为什么要写自己的？我们实际上实在官方的基础上添加了一些：request和response，handler等等
type Context struct {
	//我还把request和response封装到了这个context结构体里
	//请求和返回
	request *http.Request //这里为什么用指针:
	//https://blog.csdn.net/qq_34556414/article/details/122296968
	//如果是比较大的结构体，每次参数传递或者调用方法都要内存拷贝，
	//内存占用多，这时候可以考虑使用指针；

	//一般不会需要指向接口的指针，应该将接口作为值传递
	//因为接口低层就是个指针
	responseWriter http.ResponseWriter

	//自定义的ctx包含标准库的ctx 接口,当然你需要实现这个接口
	//如何实现:直接调用request的ctx的实现
	ctx context.Context
	//自定的handler函数,是一个slice,存的是ControllerHandler这个字段,它是一个函数签名
	handlers []ControllerHandler
	index    int // 当前请求调用到调用链的哪个节点
	// 是否超时标记位
	hasTimeout bool
	// 写保护机制sync.Mutex是一个结构体,其他的要么是接口要么是方法
	writerMux *sync.Mutex       //mutex是个结构,所以在实际实现我的ctx的时候,需要加{}
	params    map[string]string // url路由匹配的参数

	queryCache interface{}
	Request    interface{}
}

//返回一个自定义的Context,这是一个装饰器模型?还是工厂模式应该是工厂模式,没有修饰添加功能,只是返回一个Context结构
//,对传入的一个http.Request的r 进行了修饰,他被封装成一个自己上面写的Context,包含有四个结构体,还包含了一个锁.
//构造函数
func NewContext(r *http.Request, w http.ResponseWriter) *Context {

	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{},
		//这里必须要加{},why?很简单,这里需要的是一个地址,一个实际的地址,而不加{}
		//的话并只是一个实际数据结构 sync.Mutex ,而仅仅只是一个
		index: -1,
	}
}

// 调用context的pipelin上下一个handler,因为context使用了一个slice链式存储handler函数
//我们需要一个函数,来执行slice顺序存储下一个函数,也就是Next() 方法
//当然这个函数需要在入口处调用+中间件处调用
func (ctx *Context) Next() error {
	ctx.index++                        //每次调用next,指针指向下一位,
	if ctx.index < len(ctx.handlers) { //判断是否走到底,走到底就不执行了
		//是一个common ok的格式,,common中的判断实际上已经执行了这个调用中的handler方法了
		//只是handlers方法还会返回一个error,我们还能判断这个是否报错
		//所以你的handler中间件要自己调用Next()方法,这样就形成了一个pipiline调用链
		if err := ctx.handlers[ctx.index](ctx); err != nil { //这里就已经执行下一个函数了
			return err
		}
	}
	return nil
}

// #region base function,这些基本函数其实就是来自于别的包,我只是封装了下
//用的官方的锁,我们需要所,目的是为了防止重复写一个json到response里
//这些注册只是先放到这里,我们可能还不会去使用
func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

//用的传入的Request
func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

//设置hasTimeout为ture
func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

//设置hasTimeout为ture
func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

//有timeout
func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

// 设置参数
func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

// #endregio

//这个函数返回一个基本的ctx以供使用
func (ctx *Context) BaseContext() context.Context {
	//直接使用request的Context()函数,返回的是context.Context结构体
	//return ctx.request.Context()
	//这里为什么用Background?tmd竟然也是可以的
	return context.Background()
}

//这个Done函数其实request没有实现(因为request只是一个抽象类,它并没有具体实现),
//那么编译器就会去上一层去找Done方法,这叫做代理delegate,最后找到了标准库的Done()方法
//因为我们type embeding 了官方ctx,所以我们也可以直接使用代理,我们自己不实现Done,value,err,Deadlin
func (ctx *Context) Done() <-chan struct{} {
	//同样Done函数也用BaseContext()的Done()(其实是request的Done(),但是request没有这个方法,
	//就又会去标准库ctx找Done函数),就不用我们去实现了
	return ctx.BaseContext().Done()
}

//其实也是request没有实现的函数,会去找上层,也就是找标准库的Deadline()方法,
//标准库中Background和ToDo的实际上就是emptyCtx的方法
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

//同上
func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

//同上
func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

//以下两个函数在request中使用cast库重写了
// #region query url 这个函数是查找url的一些条件
//query url的几个方法都是基于Queryall方法,这个query all方法 又是调用的request.URL.Query()
//在query all的基础上 修饰,修饰到符合自身的情况
//func (ctx *Context) QueryInt(key string, def int) int {
//	params := ctx.QueryAll()
//	if vals, ok := params[key]; ok {
//		len := len(vals)
//		if len > 0 {
//			intval, err := strconv.Atoi(vals[len-1])
//			if err != nil {
//				return def
//			}
//			return intval
//		}
//	}
//	return def
//}

//func (ctx *Context) QueryString(key string, def string) string {
//	params := ctx.QueryAll()
//	if vals, ok := params[key]; ok {
//		len := len(vals)
//		if len > 0 {
//			return vals[len-1]
//		}
//	}
//	return def
//}

func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}
func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.Request != nil {
			//c.queryCache = c.Request.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}
// #end of region query url

// #region form post
//form post也基本上就是基于FormAll方法,在这个方法的基础上修修补补
//form 应该是为了返回格式
func (ctx *Context) FormInt(key string, def int) int {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			intval, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) FormString(key string, def string) string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) FormArray(key string, def []string) []string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}


// #endregion

// #region application/json post
// BindJson就是将 body 文本解析到 obj 结构体中
func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		// 读取文本
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		// 重新填充 request.Body，为后续的逻辑二次读取做准备
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		// 解析到 obj 结构体中
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}

// #endregion

// #region response

func (ctx *Context) IJson(obj interface{}) IResponse {
	//Marshal returns the JSON encoding of v,也就是
	//把想写的数据转为json,传给byt,然后byt在写入Response
	byt, err := json.Marshal(obj)
	if err != nil {
		return ctx.ISetStatus(http.StatusInternalServerError)
	}
	ctx.ISetHeader("Content-Type", "application/json")
	ctx.responseWriter.Write(byt)
	return ctx
}

//func (ctx *Context) HTML(status int, obj interface{}, template string) error {
//
//	return nil
//}

func (ctx *Context) IText(format string, values ...interface{}) IResponse {
	out := fmt.Sprintf(format, values...)
	ctx.ISetHeader("Content-Type", "application/text")
	ctx.responseWriter.Write([]byte(out))
	return ctx
}

// #endregion

// 为context设置handlers
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

func (ctx *Context) ISetHeader(key string, val string) IResponse {
	ctx.responseWriter.Header().Add(key, val)
	return ctx
}

func (ctx *Context) IXml(obj interface{}) IResponse {
	return nil
}

func (ctx *Context) IRedirect(path string) IResponse {

	return nil
}

//在Response.go有具体实现了 这里就可以删掉了,本来就是个空实现
//func (ctx *Context) Html(template string, obj interface{}) IResponse{
//	return nil
//}

func (ctx *Context) ISetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {

	return nil
}

func (ctx *Context) ISetStatus(code int) IResponse {

	return nil
}

func (ctx *Context) ISetOkStatus() IResponse {
	return nil
}

*/
