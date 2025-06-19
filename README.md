# gin-mall

> 项目来源：https://github.com/CocaineCong/gin-mall
## Todo
- [ ] 完善订单模块的状态机
- [ ] 改进服务类的单例模式
- [ ] AES 对称加密金额的密钥 key 不应当由客户端提供

## docker 部署
在项目目录下，根据 `Dockerfile` 文件构建 gin-mall 镜像：
```bash
sudo docker build -t gin-mall:v3.1 .
```

在项目目录下，加载 `compose.yaml` 配置文件启动应用:
```bash
sudo docker compose -f compose.yaml up -d
```
使用 `compose down` 指令下线应用:
```bash
sudo docker compose -f compose.yaml down
```
使用 docker 构建的容器目录树如下:
```
/app
|
|---- conf
|   |----mysql_master
|   |   |----my.cnf # 主库配置文件
|   |----mysql_master_init
|   |   |----init.sql # 主库初始化脚本
|   |----mysql_slave
|   |   |----my.cnf # 从库配置文件
|   |----mysql_slave_init
|   |   |----init.sql # 从库初始化脚本
|   |----config.go
|   |----config.ini
|---- logs
|   |---- *.log # 日志文件
|---- main
```
容器外部映射端口： \
`gin-mall`: `3000` \
`mysql-master`:`13306` \
`mysql-slave`:`13307` \
`redis`:`16379`
## MySQL 读写分离
主库名为 `mysql_master`, 从库名为 `mysql_slave`。\
指定的目录挂载和卷映射:
```yaml
volumes: # mysql_master
  - mysql_master_data:/var/lib/mysql
  - ./conf/mysql_master/my.cnf:/etc/mysql/my.cnf
  - ./conf/mysql_master_init:/docker-entrypoint-initdb.d
```
```yaml
volumes: # mysql_slave
  - mysql_slave_data:/var/lib/mysql
  - ./conf/mysql_slave/my.cnf:/etc/mysql/my.cnf
  - ./conf/mysql_slave_init:/docker-entrypoint-initdb.d
```
