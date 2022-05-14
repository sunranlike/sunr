package contract

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sunranlike/sunr/framework"
)

const RedisKey = "hade:redis"

// RedisOption 代表初始化的时候的选项
//这里指的是对Redis的配置函数，有点类似gorm的Option了
//这里传入的实redisConfig指针，就有能耐修改config
type RedisOption func(container framework.Container, config *RedisConfig) error

// RedisService 表示一个redis服务
//接口，只有一个方法，就是GetClinet方法，这个方法返回redis实例化对象
//这里这样写的原因就是很直观，某个方法里面参数可以设置为RedisService，然后传入的值就必须是实现这个接口的结构体。
type RedisService interface {
	GetClient(options ...RedisOption) (*redis.Client, error)
}

// RedisConfig 为hade定义的Redis配置结构
type RedisConfig struct {
	*redis.Options
}

// UniqKey 用来唯一标识一个RedisConfig配置
func (config *RedisConfig) UniqKey() string { //sprintf返回一个字符串而不会有任何输出。
	return fmt.Sprintf("%v_%v_%v_%v", config.Addr, config.DB, config.Username, config.Network)
}
