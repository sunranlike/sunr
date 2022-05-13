package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/sunranlike/hade/framework"
	"github.com/sunranlike/hade/framework/contract"
	"sync"
)

// HadeRedis 代表hade框架的redis实现
type HadeRedis struct {
	container framework.Container      // 服务容器
	clients   map[string]*redis.Client // key为uniqKey, value为redis.Client (连接池）

	lock *sync.RWMutex
}

// NewHadeRedis 代表实例化Client
func NewHadeRedis(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	clients := make(map[string]*redis.Client)
	lock := &sync.RWMutex{}
	return &HadeRedis{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil
}

func (app *HadeRedis) GetClient(option ...contract.RedisOption) (*redis.Client, error) {
	//从container中获取config
	config := GetBaseConfig(app.container)
	//遍历传入的option并且执行option，当然要将config传入
	for _, opt := range option {
		err := opt(app.container, config) //遍历option,并且执行返回的配置函数
		if err != nil {
			return nil, err
		}
	}
	//获取key,key应该是唯一的
	key := config.UniqKey()
	//这时候要去读map了,所以要加锁,这里锁是读写锁,读写锁读多写少
	app.lock.Lock()
	//从map容器中查询是否有服务,key就是之前的unikey
	if db, ok := app.clients[key]; ok { //如果map中存在的，直接取出并且返回
		app.lock.Unlock()
		return db, nil
	}
	app.lock.Unlock() //读取完了,就要解锁
	// 没有实例化gorm.DB，那么就要进行实例化操作
	app.lock.Lock()         //又来加锁?这里真的有必要吗?增大锁的粒度?避免性能下降.
	defer app.lock.Unlock() //取锁紧跟着解锁,是个很好的操作,当然这仅限于没有重入锁危险
	//map 中没有实例化对象，需要在map中新建一个reidis对象
	client := redis.NewClient(config.Options)
	app.clients[key] = client //存入客户端map中
	return client, nil

}
