package framework

import (
	"errors"
	"sync"
)

//https://blog.csdn.net/wuxianbing2012/article/details/122940843
//Laravel 服务容器是一个用于管理类依赖和执行依赖注入的强大工具。
//依赖注入听上去很花哨，其实质是通过构造函数或者某些情况下通过「setter」方法将类依赖注入到类中。
//简单的说服务容器就是管理类的依赖和执行依赖注入的工具，这是官方文档上说的
//我的理解更偏向于：一段生命周期所抽象的一个对象
//很难理解，打个比方，在一次请求中，你可能会用到很多服务，比如路由，队列，中间件，自定义服务等等。
//那么如何这么多的服务，如果不能妥善管理，势必会造成各种混乱，
//这是你不希望看到的，于是你想出了一个办法，用一个篮子，来乘放这些服务，然后使用的时候从篮子中拿就行了。

// Container 是一个服务容器，提供绑定服务和获取服务的功能
//我们要在framework中定义container接口和provider接口
//并且要在框架文件夹内定义一个实现container接口的结构体。
//但是provide接口不需要在框架文件夹内实现
//这样实现了核心容器只有一个实现，而你的服务可以有多种实现。
//
type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，返回 error
	Bind(provider ServiceProvider) error
	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务，
	Make(key string) (interface{}, error)

	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会 panic。
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) interface{}
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的 params 参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// HadeContainer 是服务容器的具体实现,服务可以注册入这个结构
//这个结构也会作为返回值返回给gin的engine字段的Container接口
//当然就需要这个结构体实现接口描述的方法
type HadeContainer struct {
	Container // 强制要求 HadeContainer 实现 Container 接口
	// providers 存储注册的服务提供者，key 为字符串凭证
	providers map[string]ServiceProvider
	// instance 存储具体的实例，key 为字符串凭证
	instances map[string]interface{}
	// lock 用于锁住对容器的变更操作
	lock sync.RWMutex
}

//下面对framework.container接口进行实现。

// Bind 将服务容器和关键字做了绑定,
//该函数做绑定,将服务提供者注册到provider容器之中,name作为key,ServiceProvider作为value
//然后判断是否延迟实例化,如果不是延迟实例化,注册的时候就需要实例化
//加锁--->取参数的ServiceProvider的name---->绑定到接受者的providersmap中
func (hade *HadeContainer) Bind(provider ServiceProvider) error {
	hade.lock.Lock() //加锁，并且是读写锁。
	//defer hade.lock.RUnlock()
	//这样下面的Boot不可重入锁
	key := provider.Name() //获取服务名

	hade.providers[key] = provider //name:provider存入map中
	//这里直接解锁,不然Boot调用的MustMake不会取得锁
	hade.lock.Unlock()
	//如果不需要注册时实例化,这里就结束了

	// 如果注册时就要实例化,这里开始进行实例化
	//我们的框架服务的IsDefer都是false
	if provider.IsDefer() == false { //不延迟实例化的话，那么立马实例化出来这个服务。
		if err := provider.Boot(hade); err != nil { //这个服务Bind函数中有实例化的方法
			return err
		}
		// 实例化方法,因为这里已经声明了注册时既要实现
		params := provider.Params(hade)    //参数获取
		method := provider.Register(hade)  //注册NewInstance实例化方法，除了此处还有newInstance也会使用。
		instance, err := method(params...) //调用NewInstance实例化方法
		if err != nil {                    //捕获错误
			return errors.New(err.Error())
		}
		hade.instances[key] = instance //存入实例化map
	}

	return nil
}

func (hade *HadeContainer) IsBind(key string) bool {
	return hade.findServiceProvider(key) != nil
}
func (hade *HadeContainer) findServiceProvider(key string) ServiceProvider {
	hade.lock.RLock()
	defer hade.lock.RUnlock()
	if sp, ok := hade.providers[key]; ok {
		return sp
	}
	return nil
}

// Make 方式调用内部的 make 实现
func (hade *HadeContainer) Make(key string) (interface{}, error) {
	return hade.make(key, nil, false)
}
func (hade *HadeContainer) MustMake(key string) interface{} {
	serv, err := hade.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return serv
}

// MakeNew 方式使用内部的 make 初始化
func (hade *HadeContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return hade.make(key, params, true)
}

// 真正的实例化一个服务
func (hade *HadeContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	hade.lock.RLock()
	defer hade.lock.RUnlock()
	// 查询是否已经绑定注册了这个服务提供者，如果没有注册，则返回错误，你都没有绑定注册你申请个p
	sp := hade.findServiceProvider(key) //
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew { //强制创建新实例,不用已生成的实例
		return hade.newInstance(sp, params)
	}

	// 不需要强制重新实例化，如果容器中已经实例化了，那么就直接使用容器中的实例
	//hade.instances[]就是一个存放着实例化对象的池子，你想要make一个服务，就去这个里面找你要的服务，有的话直接返回。
	if ins, ok := hade.instances[key]; ok { //common ok 语法,判断是否存在一个实例化了的服务
		return ins, nil //存在返回找个服务
	}

	// 你要的实例我们的实例容器中还未实例化，则进行一次实例化
	inst, err := hade.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}

	hade.instances[key] = inst //实例服务存入实例容器中,并将这个返回
	return inst, nil
}

//估计这里是要去使用register方法？因为我们还没有把方法注册
func (hade *HadeContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	// force new a
	if err := sp.Boot(hade); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(hade)
	}
	method := sp.Register(hade)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, err
}

//这个函数是干嘛的？我们需要将容器融入gin框架，所以在gin的Engine结构体中嵌入了一个container字段，
//他是framework.Container接口，所以我们需要一个结构体实现这个接口的方法。这就是我们的HadeContainer,上面的几个方法都是为了实现container接口
//为什么用一个函数？简洁？
func NewHadeContainer() *HadeContainer {
	return &HadeContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

// NameList 列出容器中所有服务提供者的字符串凭证
func (hade *HadeContainer) NameList() []string {
	ret := []string{}
	for _, provider := range hade.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}
