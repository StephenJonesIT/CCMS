/*
 * @File: common.redis.go
 * @Description: Defines redis information of the service
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package common

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)
 
 var (
	 redisClient *redis.Client
	 ctx         = context.Background()
 )
 
 func InitRedis(addr, password string, db int) {
	 redisClient = redis.NewClient(&redis.Options{
		 Addr:     addr,
		 Password: password,
		 DB:       db,
	 })
 }
 
 func GetCache(key string) (string, error) {
	 if redisClient == nil {
		 return "", errors.New("redis client not initialized")
	 }
	 return redisClient.Get(ctx, key).Result()
 }
 
 func SetCache(key string, value interface{}, expiration time.Duration) error {
	 if redisClient == nil {
		 return errors.New("redis client not initialized")
	 }
	 return redisClient.Set(ctx, key, value, expiration).Err()
 }

func DeleteCache(key string) error {
	ctx := context.Background()
	return redisClient.Del(ctx,key).Err()
}