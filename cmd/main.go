package main

import (
	"gin-mail/conf"
	"gin-mail/routes"
)

func main() {
	conf.Init()              // 初始化 Mysql Redis 等配置
	r := routes.NewRouter()  // 处理路由
	_ = r.Run(conf.HttpPort) // :3000
}
