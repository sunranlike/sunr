/*
这个是自己写的context,为的就是方便使用,当然自己的ctx需要实现官方包的ctx的方法
然而这个方法并不需要你手动实现,你只需要直接用*http.Request包的context
怎么用:使用BaseContext()函数,他直接
*/

package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 自定义的 Context,使用request的ctx作为基本的ctx
//为什么要写自己的？我们实际上实在官方的基础上添加了一些：request和response，handler等等
type Context struct {
	//请求和返回
	request        *http.Request //这里为什么用指针
	responseWriter http.ResponseWriter
	//自定义的ctx包含标准库的ctx
	ctx context.Context
	//自定的handler函数
	handlers []ControllerHandler //这里变成一个数组了
	index    int                 // 当前请求调用到调用链的哪个节点
	// 是否超时标记位
	hasTimeout bool
	// 写保护机制sync.Mutex是一个结构体,其他的要么是接口要么是方法
	writerMux *sync.Mutex //mutex是个结构,所以在实际实现我的ctx的时候,需要加{}

}

//返回一个自定义的Context,这是一个装饰器模型?还是工厂模式应该是工厂模式,没有修饰添加功能,只是返回一个Context结构
//,对传入的一个http.Request的r 进行了修饰,他被封装成一个自己上面写的Context,包含有四个结构体,还包含了一个锁.
func NewContext(r *http.Request, w http.ResponseWriter) *Context { //构造函数,工厂模式
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{}, //这里必须要加{},why?很简单,这里需要的是一个地址,一个实际的地址,而不加{}
		//的话并只是一个实际数据结构 sync.Mutex ,而仅仅只是一个
		index: -1,
	}
}

// 核心函数，调用context的下一个函数
func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil { //这里就已经执行下一个函数了
			return err
		}
	}
	return nil
}

// #region base function,这些基本函数其实就是来自于别的包,我只是封装了下
//用的官方的锁,我们需要所,目的是为了防止重复写一个json到response里
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

// #endregio

//BaseContext()调用的其实是requestr里面的Context()函数,这个函数返回的还是基本ctx
//不过由于返回的是request里面的ctx,所有我们的baseCtx就是有一些request的方法
func (ctx *Context) BaseContext() context.Context {
	//直接使用request的Ctx,导入这个标准包,我就不用再去写新的实现了
	return ctx.request.Context()
}

//这个Done函数其实request没有实现,那么编译器就会去上一层去找Done,这叫做代理delegate
func (ctx *Context) Done() <-chan struct{} {
	//同样Done函数也用BaseContext()的Done()(其实是request的Done(),但是request没有这个方法,
	//就又会去标准库ctx找Done函数),就不用我们去实现了
	return ctx.BaseContext().Done()
}

//其实也是request没有实现的函数,会去找上层,也就是找标准库的Deadline()方法
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

// #region query url 这个函数是查找url的一些条件
//query url的几个方法都是基于Queryall方法,这个query all方法 又是调用的request.URL.Query()
//在query all的基础上 修饰,修饰到符合自身的情况
func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
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

func (ctx *Context) QueryString(key string, def string) string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

//QueryAll() 方法调用了request的URL.QUERY()的方法
//
func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.URL.Query())
	}
	return map[string][]string{}
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

//

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}

// #endregion

// #region application/json post

func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

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
//status是返回的http状态码, obj是写入的字符串？也不一定是字符串
func (ctx *Context) Json(status int, obj interface{}) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj) //Marshal returns the JSON encoding of v,也就是八字
	//把想写的数据转为json,传给byt,然后byt在写入Response
	if err != nil {
		ctx.responseWriter.WriteHeader(500)
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (ctx *Context) Text(status int, obj string) error {
	return nil
}

// #endregion

// 为context设置handlers
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}
