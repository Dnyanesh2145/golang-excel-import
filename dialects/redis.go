package dialects

import (
	"context"
	"fmt"
	"golang-excel-import/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisc struct {
	rdb *redis.Client
}

var RedisClient = &redisc{}
var ctx = context.Background()

func (r *redisc) Connect() {
	if options, err := redis.ParseURL(config.GetLocalEnv("REDIS_URI")); err != nil {
		fmt.Println("Unable to Parse Redis URI: ", err)
	} else {
		r.rdb = redis.NewClient(options)
	}
	isConnected := r.rdb.Ping(ctx).String()
	if isConnected == "ping: PONG" {
		fmt.Println("Connected to redis successfully")
	} else {
		fmt.Println("Unable to connect Redis: ", isConnected)
	}
}

// Get retrieves the value stored at the specified key in Redis.
func (r *redisc) Get(key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

// method which will set keys with TTL
func (r *redisc) SetE(key string, value string, ttl time.Duration) error {
	err := r.rdb.SetEx(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}
