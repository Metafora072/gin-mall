package model

// BasePage 描述分页
// PageNum: 页的数量
// PageSize: 每页的条目数量
type BasePage struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}
