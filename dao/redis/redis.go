package redis

import (
	"fmt"
	"web_app/settings"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

func Init(cfg *settings.RedisConfig) (err error) {
	//rdb = redis.NewClient(&redis.Options{
	//	Addr: fmt.Sprintf("%s:%d",
	//		viper.GetString("redis.host"),
	//		viper.GetInt("redis.port")),
	//	Password: viper.GetString("redis.password"),
	//	DB:       viper.GetInt("redis.db"),
	//	PoolSize: viper.GetInt("redis.pool_size"),
	//})
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	_, err = client.Ping().Result()
	return
}
func Close() {
	_ = client.Close()
}
