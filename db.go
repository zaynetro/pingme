package main

import (
	"github.com/garyburd/redigo/redis"
)

func dial(network string, address string, password string) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}

	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}

	return c, err
}

func dialLocal() redis.Conn {
	c, err := dial("tcp", CONFIG.Redis.Host, CONFIG.Redis.Password)
	if err != nil {
		panic(err)
	}
	return c
}

func getString(command string, args ...interface{}) (string, error) {
	conn := dialLocal()
	defer conn.Close()
	return redis.String(conn.Do(command, args...))
}

func exec(command string, args ...interface{}) {
	conn := dialLocal()
	defer conn.Close()
	conn.Do(command, args...)
}

func exists(key string) bool {
	conn := dialLocal()
	defer conn.Close()
	val, _ := redis.Bool(conn.Do("EXISTS", key))
	return val
}
