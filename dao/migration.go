package dao

import "gin-mail/model"

// migrate 将 struct 结构体迁移成 MySQL 表结构
// 使用 GORM 进行数据库迁移，确保 UserModel 和 TaskModel 这两个数据模型的表在数据库中存在，并且自动更新或创建数据库表结构。
func migration() {
	// Set 方法是 GORM 提供的设置方法，用来设置数据库连接的选项。
	// gorm:table_options 是 GORM 的配置项，用来设置表的特定选项。在这个例子中，设置了表的字符集为 utf8mb4，这是支持多语言字符集（包括 Emoji）的字符集。
	// AutoMigrate 是 GORM 中用于 自动迁移 数据库表结构的方法。
	// 它会根据你定义的模型来检查数据库中是否存在对应的表结构.
	// 如果不存在，GORM 会自动创建；如果存在，GORM 会检查表结构和模型定义的差异，必要时会自动更新表结构。
	err := _db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
		&model.Address{},
		&model.Admin{},
		&model.Carousel{},
		&model.Cart{},
		&model.Category{},
		&model.Favorite{},
		&model.Notice{},
		&model.Order{},
		&model.Product{},
		&model.ProductImg{},
		&model.User{},
	)

	if err != nil {
		panic(err)
	}
}
