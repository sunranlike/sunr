package demo

// Demo 服务的 key
const Key = "hade:demo"

// Demo 服务的接口
type Service interface {
	GetFoo() Foo
}

// Demo 服务接口定义的一个数据结构
type Foo struct {
	Name string
}

//provider文件实现ServiceProvider接口，让服务能够注册入框架
//Contract文件定义我的服务具体是干嘛的，有哪些方法个接口，比如GetFoo（）方法
//同时Contract还能定义一些provider用的参数比如说Key和一些其他结构体

//service文件根据contract实现具体协议，同时实现注册的返回函数，一般来说provider文件中并不直接写注册函数
//注册函数也在service文件中写。
