<<<<<<< HEAD

=======
>>>>>>> 8e475c8 (Initial commit)
package framework

import "net/http"

// 框架核心结构
type Core struct {
}

// 初始化框架核心结构
func NewCore() *Core {
	return &Core{}
}

// 框架核心结构实现Handler接口，如果你不自己实现就会调用默认Handler，
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// TODO
}
<<<<<<< HEAD
=======

>>>>>>> 8e475c8 (Initial commit)
//在源码中Handler 接口实际上只有一个函数 就是ServerHTTP方法，所以我们自己写一个ServeHttp就代表我们
//使用自己的Handler，（实际上是一个ServerHttp方法）
//type Handler interface {
//	ServeHTTP(ResponseWriter, *Request)
<<<<<<< HEAD
//}
=======
//}
>>>>>>> 8e475c8 (Initial commit)
