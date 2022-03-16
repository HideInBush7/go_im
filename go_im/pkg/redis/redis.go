package redis

import (
	"time"

	"github.com/HideInBush7/go_im/pkg/config"
	"github.com/gomodule/redigo/redis"
)

type redisConf struct {
	Address     string `json:"address"`
	MaxIdle     int    `json:"maxIdleConns"`
	IdleTimeout int    `json:"idleTimeout"`
}

var redisPool *redis.Pool

func init() {
	// 从配置文件(viper)中key为'redis'获取配置
	var conf = redisConf{}
	config.Register(`redis`, &conf)

	if conf.MaxIdle <= 0 {
		conf.MaxIdle = 3
	}
	if conf.IdleTimeout <= 0 {
		conf.IdleTimeout = 240
	}
	redisPool = &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			return redis.Dial(`tcp`, conf.Address)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do(`ping`)
			return err
		},
	}
}

// 获取redis连接实例
func GetInstance() redis.Conn {
	return redisPool.Get()
}
