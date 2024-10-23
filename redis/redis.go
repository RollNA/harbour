package redis

import (
	"context"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var Cluster *redis.ClusterClient
var Client *redis.Client

type RedisConf struct {
	Addr     string `json:"addr" mapstructure:"addr"`
	Password string `json:"password" mapstructure:"password"`
}

func InitCluster(conf RedisConf) {
	newClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{conf.Addr},
		Password: conf.Password,
	})

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(newClient); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(newClient); err != nil {
		panic(err)
	}
	pong := newClient.Ping(context.Background()).Val()
	if pong != "PONG" {
		panic("redis " + pong)
	}
	Cluster = newClient
}

func Init(conf RedisConf) {
	newClient := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
	})

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(newClient); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(newClient); err != nil {
		panic(err)
	}
	pong, err := newClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	if pong != "PONG" {
		panic("redis ping:" + pong)
	}
	Client = newClient
}
