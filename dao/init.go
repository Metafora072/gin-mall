package dao

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"time"
)

var _db *gorm.DB

func Database(connRead, connWrite string) {
	// 根据当前 Gin 的运行模式（gin.Mode()）来设置 ORM 的日志级别
	// gin.Mode() 返回当前 Gin 框架的运行模式，它可以是以下几种之一：
	// "debug"：开发模式，通常会输出详细的调试信息。
	// "release"：生产模式，通常只输出错误信息。
	// "test"：测试模式，通常会禁止日志输出。
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connRead, // DSN data source name
		DefaultStringSize:         256,      // string 类型字段的默认长度
		DisableDatetimePrecision:  true,     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,    // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: ormLogger, // 打印日志
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表明不加s, user ->× users
		},
	})

	if err != nil {
		panic(err)
	}

	// 从 GORM 实例 db 中获取底层的 *sql.DB 对象（标准的 Go 数据库连接对象）。GORM 在内部使用 *sql.DB 来处理实际的数据库连接和操作。
	sqlDB, _ := db.DB()
	// 设置数据库连接池中 最大空闲连接数。空闲连接是指在连接池中处于未使用状态的连接。
	// 这里设置了最大空闲连接数为 20，意味着连接池最多可以保持 20 个空闲连接，如果空闲连接数超过这个值，连接会被关闭。
	sqlDB.SetMaxIdleConns(20)
	// 设置数据库连接池中 最大打开连接数。这个限制了同时与数据库建立的最大连接数。
	// 这里设置了最大连接数为 100，意味着最多允许 100 个并发的数据库连接被打开。如果超过这个限制，新的数据库请求将会被阻塞，直到有空闲连接可用。
	sqlDB.SetMaxOpenConns(100)
	// 设置数据库连接的 最大生命周期。当连接池中的连接超过这个生命周期后，它们会被关闭并且重新创建。
	// 这里设置最大连接生命周期为 30 秒 (time.Second * 30)，意味着连接池中的每个连接最多存在 30 秒，超过时间后就会被关闭并重建。这个设置有助于避免长时间占用的数据库连接可能引起的资源泄漏或数据库连接不稳定问题。
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	_db = db

	// 主从配置
	_ = _db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(connWrite)},                      // 写操作
		Replicas: []gorm.Dialector{mysql.Open(connRead), mysql.Open(connRead)}, // 读操作
		Policy:   dbresolver.RandomPolicy{},
	}))

	migration()
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
