package cache

import (
	"context"
	"gin-mail/pkg/utils"
	"github.com/go-redis/redis/v8"
	"gopkg.in/ini.v1"
	"strconv"
)

var RedisClient *redis.Client

var (
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

func RedisInit() {
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		utils.LogrusObj.Infoln("Redis ini Load:", err)
		panic(err)
	}
	LoadRedis(file)

	db, _ := strconv.ParseInt(RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPw,
		DB:       int(db),
	})

	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		utils.LogrusObj.Infoln("redis ping:", err)
		panic(err)
	}
	RedisClient = client
}

func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}
