package main

import (
	"gin-mail/cache"
	"gin-mail/conf"
	"gin-mail/pkg/utils"
	"gin-mail/routes"
)

func main() {
	conf.Init()              // 初始化 Mysql, path, email, service 等配置
	cache.RedisInit()        // 初始化 Redis 配置
	utils.InitLog()          // 初始化日志对象
	r := routes.NewRouter()  // 处理路由
	_ = r.Run(conf.HttpPort) // :3000
}
