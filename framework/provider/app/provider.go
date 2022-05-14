package app

import (
	"github.com/sunranlike/sunr/framework"
	"github.com/sunranlike/sunr/framework/contract"
)

// HadeAppProvider 提供App的具体实现方法
//这个结构提供一个基础的目录服务,属于框架级别服务
type HadeAppProvider struct {
	BaseFolder string
}

// Register 注册HadeApp方法
func (h *HadeAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewHadeApp //实例化注册方法需要在service中具体实现.
}

// Boot 启动调用
func (h *HadeAppProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化。如果返回false就是在绑定的时候就实例化
func (h *HadeAppProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (h *HadeAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, h.BaseFolder}
}

// Name 获取字符串凭证
func (h *HadeAppProvider) Name() string {
	return contract.AppKey
}
