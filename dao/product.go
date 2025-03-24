package dao

import (
	"context"
	"gin-mail/model"
	"gorm.io/gorm"
)

type ProductDao struct {
	*gorm.DB
}

func NewProductDao(ctx context.Context) *ProductDao {
	return &ProductDao{NewDBClient(ctx)}
}

func NewProductDaoByDB(db *gorm.DB) *ProductDao {
	return &ProductDao{db}
}

func (dao *ProductDao) CreateProduct(product *model.Product) (err error) {
	return dao.DB.Model(&model.Product{}).Create(&product).Error
}

// CountProductByCondition 查找满足指定条件 condition 的商品数量有多少个
func (dao *ProductDao) CountProductByCondition(condition map[string]interface{}) (total int64, err error) {
	err = dao.DB.Model(&model.Product{}).Where(condition).Count(&total).Error
	return
}

// ListProductByCondition 根据指定条件 condition 和分页模式 page (pageNum 表示展示第几页，pageSize 表示一页展示几个商品) 展示商品列表
func (dao *ProductDao) ListProductByCondition(condition map[string]interface{}, page model.BasePage) (products []*model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).Where(condition).Offset((page.PageNum - 1) * page.PageSize).Limit(page.PageSize).Find(&products).Error
	return
}

// SearchProduct 根据搜索信息 info 搜索符合的商品列表并以指定的分页模式 page 来展示
// count: 符合条件的商品总数
func (dao *ProductDao) SearchProduct(info string, page model.BasePage) (products []*model.Product, count int64, err error) {
	err = dao.DB.Model(&model.Product{}).Where("title LIKE ? OR info LIKE ?", "%"+info+"%", "%"+info+"%").Count(&count).Error
	if err != nil {
		return
	}

	// TODO: 会导致全表扫描，性能问题？ 可以考虑：为 title 和 info 字段创建 全文索引（Fulltext Index）；使用 MATCH AGAINST（如 MySQL 的全文搜索）提升查询性能。
	// "%"+info+"%": SQL 的 LIKE 语法，% 是通配符，表示匹配任意字符（如 %test% 匹配包含 "test" 的字符串）。
	err = dao.DB.Model(&model.Product{}).Where("title LIKE ? OR info LIKE ?", "%"+info+"%", "%"+info+"%").Offset((page.PageNum - 1) * page.PageSize).Limit(page.PageSize).Find(&products).Error
	return
}
