package main

import (
	"github.com/garyburd/redigo/redis"
)

func dial(network string, address string) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
	// defer c.Close()
	return c, err
}

func dialLocal() redis.Conn {
	c, err := dial("tcp", ":6379")
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
