package db

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"phyFit/mqtt/config"
)

var redisPool *redis.Pool

func init() {
	RedisPollInit()
}

func RedisPollInit() *redis.Pool {
	//var redisPool *redis.Pool
	redisPool = &redis.Pool{
		MaxIdle:   3,
		MaxActive: 0,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			// c, err := redis.Dial("tcp", "127.0.0.1:6379")
			c, err := redis.Dial("tcp", config.RedisConfig.Addr)
			if err != nil {
				return nil, err
			}
			// if _, err = c.Do("AUTH", "123456"); err != nil {
			if _, err = c.Do("AUTH", config.RedisConfig.Password); err != nil {
				c.Close()
				fmt.Println("redis password failed", err)
				return c, err
			}
			return c, err
		},
	}
	return redisPool
}

func AddCache(key, value string) error {

	conn := redisPool.Get()
	fmt.Printf("Connent Info --- %+v", conn)
	defer conn.Close()

	if _, err := conn.Do("SET", key, value); err != nil {
		fmt.Println("error -- ", err)
		return err
	}
	return nil
}

func GetKey(str string) (str1 string, err error) {
	conn := redisPool.Get()
	s, err := redis.String(conn.Do("GET", str))
	if err != nil {
		fmt.Printf("err : %+v\n", err)
		return "", err
	}
	defer conn.Close()

	return s, nil
}

// Delete Cache
func DeleteCache(key string) error {
	conn := redisPool.Get()
	defer conn.Close()
	if _, err := conn.Do("DEL", key); err != nil {
		return err
	}
	return nil
}
