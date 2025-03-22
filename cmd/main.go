package main

import (
	"gin-mail/conf"
	"gin-mail/pkg/utils"
	"gin-mail/routes"
)

func main() {
	conf.Init()              // 初始化 Mysql Redis 等配置
	utils.InitLog()          // 初始化日志对象
	r := routes.NewRouter()  // 处理路由
	_ = r.Run(conf.HttpPort) // :3000
}
