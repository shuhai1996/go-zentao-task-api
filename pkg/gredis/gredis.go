package gredis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"go-zentao-task/pkg/config"
	"log"
	"time"
)

var RedisPool *redis.Pool
var RedisLogPool *redis.Pool

type redisConfig struct {
	host     string
	port     string
	password string
}

func Setup() {
	RedisPool = redisConnectPool("")
	RedisLogPool = redisConnectPool(".log")
}

func getRedisConfig(store string) (conf *redisConfig, err error) {
	conf = &redisConfig{
		host:     config.Get("redis.host" + store),
		port:     config.Get("redis.port" + store),
		password: config.Get("redis.password" + store),
	}
	if conf.host == "" {
		err = errors.New("redis.host" + store + "不能为空")
		return
	}
	if conf.port == "" {
		err = errors.New("redis.port" + store + "不能为空")
		return
	}
	return
}

func newRedisPool(server, password string) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     50,
		MaxActive:   100,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
	}

	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		log.Fatalln(err)
	}
	return pool
}

func redisConnectPool(store string) *redis.Pool {
	conf, err := getRedisConfig(store)
	if err != nil {
		log.Fatalln(err)
	}
	return newRedisPool(conf.host+":"+conf.port, conf.password)
}
