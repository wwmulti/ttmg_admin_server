package utils

import (
	"context"
	"encoding/json"
	"sync"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/redis/go-redis/v9"
)

var (
	redisOnce   sync.Once
	redisClient *redis.Client
)

// 获取 Redis 客户端（单例）
func GetRedisClient() *redis.Client {
	redisOnce.Do(func() {
		redisconn, err := beego.AppConfig.String("redisconn")
		if redisconn == "" || nil != err {
			panic("redisconn配置没有设置！")
		}
		
		type STConfig struct {
			Conn     string
			Password string
			DbNum    int
			PoolSize int
		}
		var cf STConfig

		err = json.Unmarshal([]byte(redisconn), &cf)
		if err != nil {
			panic("Redis配置JSON解析失败: " + err.Error())
		}

		redisClient = redis.NewClient(&redis.Options{
			Addr:     cf.Conn,
			Password: cf.Password,
			DB:       cf.DbNum,
			PoolSize: cf.PoolSize,
		})
		
		_, err = redisClient.Ping(context.Background()).Result()
		if err != nil {
			panic("Redis connect ping failed, err:" + err.Error())
		}
	})	
	return redisClient
}
