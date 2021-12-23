package service

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"go-zentao-task/pkg/db"
	"go-zentao-task/pkg/gredis"
)

func PingMysql() error {
	return db.Orm.DB().Ping()
}

func PingRedis() error {
	conn := gredis.RedisPool.Get()
	defer conn.Close()

	pong, err := redis.String(conn.Do("ping"))
	if err != nil {
		return err
	}
	if pong != "PONG" {
		return errors.New("redis ping error: " + pong)
	}

	return nil
}

func PingLogRedis() error {
	conn := gredis.RedisLogPool.Get()
	defer conn.Close()

	pong, err := redis.String(conn.Do("ping"))
	if err != nil {
		return err
	}
	if pong != "PONG" {
		return errors.New("log redis ping error: " + pong)
	}

	return nil
}
