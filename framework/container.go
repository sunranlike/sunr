package framework

import (
	"errors"
	"sync"
)

// Container 是一个服务容器，提供绑定服务和获取服务的功能
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
type HadeContainer struct {
	Container // 强制要求 HadeContainer 实现 Container 接口
	// providers 存储注册的服务提供者，key 为字符串凭证
	providers map[string]ServiceProvider
	// instance 存储具体的实例，key 为字符串凭证
	instances map[string]interface{}
	// lock 用于锁住对容器的变更操作
	lock sync.RWMutex
}

// Bind 将服务容器和关键字做了绑定,
//该函数做绑定,将服务提供者注册到provider容器之中,name作为key,ServiceProvider作为value
//然后判断是否延迟实例化,如果不是延迟实例化,注册的时候就需要实例化
//加锁--->取参数的ServiceProvider的name---->绑定到接受者的providersmap中
func (hade *HadeContainer) Bind(provider ServiceProvider) error {
	hade.lock.Lock() //加锁，并且是读写锁。
	defer hade.lock.Unlock()
	key := provider.Name() //获取服务名

	hade.providers[key] = provider //name:provider存入map中

	//如果不需要注册时实例化,这里就结束了

	// 如果注册时就要实例化,这里开始进行实例化
	if provider.IsDefer() == false { //如果这个server就是必须注册就要实例化的,那么boot中就要有他的实例化方法
		if err := provider.Boot(hade); err != nil { //这个服务Bind函数中有实例化的方法
			return err
		}
		// 实例化方法,因为这里已经声明了注册时既要实现
		params := provider.Params(hade)    //参数获取
		method := provider.Register(hade)  //注册NewInstance实例化方法
		instance, err := method(params...) //调用NewInstance实例化方法
		if err != nil {                    //捕获错误
			return errors.New(err.Error())
		}
		hade.instances[key] = instance //存入实例化map
	}
	return nil
}

// Make 方式调用内部的 make 实现
func (hade *HadeContainer) Make(key string) (interface{}, error) {
	return hade.make(key, nil, false)
}

// MakeNew 方式使用内部的 make 初始化
func (hade *HadeContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return hade.make(key, params, true)
}

// 真正的实例化一个服务
func (hade *HadeContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	hade.lock.RLock()
	defer hade.lock.RUnlock()
	// 查询是否已经注册了这个服务提供者，如果没有注册，则返回错误
	sp := hade.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew { //强制创建新实例,不用已生成的实例
		return hade.newInstance(sp, params)
	}

	// 不需要强制重新实例化，如果容器中已经实例化了，那么就直接使用容器中的实例
	if ins, ok := hade.instances[key]; ok { //common ok 语法,判断是否存在一个实例化了的服务
		return ins, nil //存在返回找个服务
	}

	// 容器中还未实例化，则进行一次实例化
	inst, err := hade.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}

	hade.instances[key] = inst //实例服务存入实例服务map中,并返回
	return inst, nil
}

func (hade *HadeContainer) findServiceProvider(key string) interface{} {

	return nil
}

func (hade *HadeContainer) newInstance(sp interface{}, params []interface{}) (interface{}, error) {
	return nil, nil
}

func NewHadeContainer() *HadeContainer {
	return &HadeContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}
