package env

import (
	"github.com/sunranlike/hade/framework"
	"github.com/sunranlike/hade/framework/contract"
)

type HadeEnvProvider struct {
	Folder string
}

// Register registe a new function for make a services instance
func (provider *HadeEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewHadeEnv
}

// Boot will called when the services instantiate
//在根据文章进行编码后, 运行程序的时候发现程序会一直卡在获取锁的地方。
//调试了很久才发现是因为 EnvProvider 的Boot()中调用了MustMake, 会再一次获取锁, 会导致死锁。
//死锁的原因是: main.go 中的Bind() 会先获取锁, 使用的是defer 释放锁, 在Bind() 中由于会调用Boot(),
//而go不支持重入锁
//EnvProvider的Boot()中也会去获取锁, 导致再次获取锁时会失败，因此会卡住.
//解决方法是: 将Bind()中的锁释放改为 hade.lock.Unlock() 直接释放，尽量让锁的占用时间最小。
func (provider *HadeEnvProvider) Boot(c framework.Container) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

// IsDefer define whether the services instantiate when first make or register
func (provider *HadeEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *HadeEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.Folder}
}

/// Name define the name for this services
func (provider *HadeEnvProvider) Name() string {
	return contract.EnvKey
}
