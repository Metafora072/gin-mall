package conf

import (
	"gin-mail/dao"
	"gin-mail/pkg/utils"
	"gopkg.in/ini.v1"
	"strings"
)

var (
	AppMode  string
	HttpPort string

	DbMaster         string
	DbHostMaster     string
	DbPortMaster     string
	DbUserMaster     string
	DbPasswordMaster string
	DbNameMaster     string

	DbSlave         string
	DbHostSlave     string
	DbPortSlave     string
	DbUserSlave     string
	DbPasswordSlave string
	DbNameSlave     string

	ValidEmail string
	SmtpHost   string
	SmtpEmail  string
	SmtpPass   string

	Host        string
	ProductPath string
	AvatarPath  string
)

func Init() {
	// 读取本地环境变量
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		panic(err)
	}
	LoadServer(file)
	LoadMysql(file)
	LoadEmail(file)
	LoadPhotoPath(file)

	// mysql 读
	slaveDsn := strings.Join([]string{DbUserSlave, ":", DbPasswordSlave, "@tcp(", DbHostSlave, ":", DbPortSlave, ")/", DbNameSlave, "?charset=utf8mb4&parseTime=True"}, "")
	// mysql 写
	masterDsn := strings.Join([]string{DbUserMaster, ":", DbPasswordMaster, "@tcp(", DbHostMaster, ":", DbPortMaster, ")/", DbNameMaster, "?charset=utf8mb4&parseTime=True"}, "")

	utils.LogrusObj.Infoln("slaveDsn:", slaveDsn)
	utils.LogrusObj.Infoln("masterDsn:", masterDsn)

	dao.Database(slaveDsn, masterDsn)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMysql(file *ini.File) {
	DbMaster = file.Section("mysql_master").Key("Db").String()
	DbHostMaster = file.Section("mysql_master").Key("DbHost").String()
	DbPortMaster = file.Section("mysql_master").Key("DbPort").String()
	DbUserMaster = file.Section("mysql_master").Key("DbUser").String()
	DbPasswordMaster = file.Section("mysql_master").Key("DbPassword").String()
	DbNameMaster = file.Section("mysql_master").Key("DbName").String()

	DbSlave = file.Section("mysql_slave").Key("Db").String()
	DbHostSlave = file.Section("mysql_slave").Key("DbHost").String()
	DbPortSlave = file.Section("mysql_slave").Key("DbPort").String()
	DbUserSlave = file.Section("mysql_slave").Key("DbUser").String()
	DbPasswordSlave = file.Section("mysql_slave").Key("DbPassword").String()
	DbNameSlave = file.Section("mysql_slave").Key("DbName").String()
}

/* LoadRedis 在 cache 模块再导入，否则会造成循环引用
func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}
*/

func LoadEmail(file *ini.File) {
	ValidEmail = file.Section("email").Key("ValidEmail").String()
	SmtpHost = file.Section("email").Key("SmtpHost").String()
	SmtpEmail = file.Section("email").Key("SmtpEmail").String()
	SmtpPass = file.Section("email").Key("SmtpPass").String()
}

func LoadPhotoPath(file *ini.File) {
	Host = file.Section("path").Key("Host").String()
	ProductPath = file.Section("path").Key("ProductPath").String()
	AvatarPath = file.Section("path").Key("AvatarPath").String()
}
