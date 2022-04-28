package demo

import (
	"fmt"
	"github.com/sunranlike/hade/framework"
)

// DemoService 具体的接口实例，这里描述服务真正的工作内容，
type DemoService struct {
	// 实现接口
	Service //显式嵌入，这样使得结构体必须实现该结构，算是增加耦合了，但是更加清晰

	// 参数
	c framework.Container
}

// 实现接口
func (s *DemoService) GetFoo() Foo {
	return Foo{
		Name: "i am foo",
	}
}

// 初始化实例的方法
func NewDemoService(params ...interface{}) (interface{}, error) {
	// 这里需要将参数展开
	c := params[0].(framework.Container)

	fmt.Println("new demo service")
	// 返回实例，有一个这个实例
	return &DemoService{c: c}, nil
}
