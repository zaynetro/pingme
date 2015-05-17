package main

import (
	"net/url"
	"os"
	"reflect"
)

type RedisConfig struct {
	Host     string `default:"redis://localhost:6379"`
	Password string `default:""`
}

type ServerConfig struct {
	Port   string `default:"3000"`
	Secret string `default:"secret"`
}

type Config struct {
	Redis  RedisConfig
	Server ServerConfig
}

// Session expiration (from github.com/boj/redistore)
const sessionExpire = 86400 * 30

var CONFIG = confApp()

func confApp() Config {
	_redis := confRedis(os.Getenv("REDIS_URL"))
	_server := confServer(os.Getenv("PORT"))

	return Config{
		Redis:  _redis,
		Server: _server,
	}
}

func confRedis(connUrl string) RedisConfig {
	_redis := RedisConfig{}
	typ := reflect.TypeOf(_redis)

	if connUrl == "" {
		h, _ := typ.FieldByName("Host")
		_redis.Host = h.Tag.Get("default")

		p, _ := typ.FieldByName("Password")
		_redis.Password = p.Tag.Get("default")

		return _redis
	}

	redisURL, err := url.Parse(connUrl)
	if err != nil {
		panic(err)
	}

	auth := ""

	if redisURL.User != nil {
		if password, ok := redisURL.User.Password(); ok {
			auth = password
		}
	}

	return RedisConfig{
		Host:     redisURL.Host,
		Password: auth,
	}
}

func confServer(port string) ServerConfig {
	_conf := ServerConfig{
		Secret: "learngo",
	}
	typ := reflect.TypeOf(_conf)

	if port == "" {
		p, _ := typ.FieldByName("Port")
		_conf.Port = p.Tag.Get("default")

		return _conf
	}

	_conf.Port = port

	return _conf
}
