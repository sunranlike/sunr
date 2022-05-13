package orm

import (
	"github.com/sunranlike/hade/framework"
	"github.com/sunranlike/hade/framework/contract"
)

//这个文件是干嘛的？其实就是思考了Gorm的Option结构，我们也设置一个DBOption
//这个DBOption作为参数传入的时候就会被遍历，然后会执行一些初始化工作。
//比如说下面的几个操作

// WithDryRun 设置空跑模式
func WithDryRun() contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		config.DryRun = true
		return nil
	}
}

// WithConfigPath 加载配置文件地址
func WithConfigPath(configPath string) contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		// 加载configPath配置路径
		if err := configService.Load(configPath, config); err != nil {
			return err
		}
		return nil
	}
}
