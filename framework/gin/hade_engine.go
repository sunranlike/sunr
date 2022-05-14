package gin

import "github.com/sunranlike/sunr/framework"

func (engine *Engine) SetContainer(container framework.Container) {
	engine.container = container
}

// GetContainer 从Engine中获取container
func (engine *Engine) GetContainer() framework.Container {
	return engine.container
}

// Bind engine实现container的绑定封装
//？这是个什么逻辑？我们在实现它的时候在调用它？
//其实engine没有实现Bind方法，他只是调用了container字段的Bind方法，
//其实也是合法的，因为规定container字段的value是一个实现了framework.Container的字段，
//那么container字段必定有一个Bind方法
//实际上我们在调用gin.New的时候，给这个container字段返回的是hadecontainer，这个歌结构体是实现了
//framerwork.Container接口的，也就是他又bind和make类方法，自然就可以用
func (engine *Engine) Bind(provider framework.ServiceProvider) error {
	return engine.container.Bind(provider)
}

// IsBind 关键字凭证是否已经绑定服务提供者
func (engine *Engine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}
