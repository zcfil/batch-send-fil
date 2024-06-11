package conf

import (
	"github.com/go-redis/redis/v7"
	"time"
)

//token在缓存中的有效时间
const (
	//api_key名称
	API_KEY string = "token"
	//缓存中token名称
	TOKEN_PREFIX string = "user_token:prefix_token_"
	//token有效期一周
	TOKEN_EFFECT_TIME int64 = 7 * 24 * 60 * 60
	//miner key名称
	MINERS_TOKEN string = "miners_token:miner_token_"
	//miner key有效期
	MINERS_EFFECT_TIME time.Duration = time.Hour * 24
	//user roles name
	ROLES_PREFIX string = "user_token:role_token_"
)

func NewRedisClient(r *Redis) func() (*redis.Client, error) {

	return func() (*redis.Client, error) {

		client := redis.NewClient(&redis.Options{
			//Network:            "",
			Addr: r.Host,
			//Dialer:             nil,
			//OnConnect:          nil,
			//Password: "kevin@163.com",
			DB: 0,
			//MaxRetries:         0,
			//MinRetryBackoff:    0,
			//MaxRetryBackoff:    0,
			//DialTimeout:        0,
			//ReadTimeout:        0,
			//WriteTimeout:       0,
			//PoolSize:           0,
			//MinIdleConns:       0,
			//MaxConnAge:         0,
			//PoolTimeout:        0,
			//IdleTimeout:        0,
			//IdleCheckFrequency: 0,
			//TLSConfig:          nil,
		})
		return client, nil
	}
}
