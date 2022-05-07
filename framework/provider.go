package framework

type NewInstance func(...interface{}) (interface{}, error)

// ServiceProvider 定义一个服务提供者需要实现的接口
//这里只需要定义接口，不需要具体的结构体去实现，
//这就是服务容器——服务提供者模型的好处
//我们只需要实现容器接口，服务提供者可以在他处实现，实现了解绑
//比如说你想要实现一个缓存服务，你只需要按照服务提供者接口去实现即可
//扩展性好且不会侵入
type ServiceProvider interface {
	// Register 在服务容器中注册了一个实例化服务的方法，是否在注册的时候就实例化这个服务，需要参考 IsDefer 接口。
	//
	Register(Container) NewInstance
	// Boot 在调用实例化服务的时候会调用，可以把一些准备工作：基础配置，初始化参数的操作放在这个里面。
	// 如果 Boot 返回 error，整个服务实例化就会实例化失败，返回错误
	Boot(Container) error
	// IsDefer 决定是否在注册的时候实例化这个服务，如果不是注册的时候实例化，那就是在第一次 make 的时候进行实例化操作
	// false 表示不需要延迟实例化，在Main中bind的时候就就要实例化。true 表示延迟实例化，一般不会有
	IsDefer() bool
	// Params params 定义传递给 NewInstance 的参数，可以自定义多个，建议将 container 作为第一个参数
	Params(Container) []interface{}
	// Name 代表了这个服务提供者的凭证
	Name() string
}
