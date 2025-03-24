package service

import (
	"context"
	"gin-mail/dao"
	"gin-mail/model"
	"gin-mail/pkg/e"
	"gin-mail/pkg/utils"
	"gin-mail/serializer"
	"mime/multipart"
	"strconv"
	"sync"
)

type ProductService struct {
	Id            uint   `json:"id" form:"id"`
	Name          string `json:"name" form:"name"`
	CategoryId    uint   `json:"category_id" form:"category_id"`
	Title         string `json:"title" form:"title"`
	Info          string `json:"info" form:"info"`
	ImgPath       string `json:"img_path" form:"img_path"`
	Price         string `json:"price" form:"price"`
	DiscountPrice string `json:"discount_price" form:"discount_price"`
	OnSale        bool   `json:"on_sale" form:"on_sale"`
	Num           int    `json:"num" form:"num"`
	model.BasePage
}

// Create 创建商品
// uid: 商品商家id
// files: 上传的多张图片文件，multipart.FileHeader 是 Gin 框架用于处理上传文件的类型
func (service *ProductService) Create(ctx context.Context, uid uint, files []*multipart.FileHeader) serializer.Response {
	var boss *model.User
	var err error
	code := e.Success

	userDao := dao.NewUserDao(ctx)
	boss, _ = userDao.GetUserById(uid)

	// 以第一张作为封面图
	tmp, _ := files[0].Open()
	// 将图片上传至本地静态目录并返回路径
	path, err := UploadProductToLocalStatic(tmp, uid, service.Name)
	if err != nil {
		code = e.ErrorProductImgUpload
		utils.LogrusObj.Infoln("ProductService func Create:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 构造 Product 数据对象, OnSale 默认为 true，表示该商品上架销售
	// 注意：product.Id 没有显式指定，在 product 指针被写入到数据库后，数据库会自增 ID 并填充回 product.Id
	product := &model.Product{
		Name:          service.Name,
		CategoryId:    service.CategoryId,
		Title:         service.Title,
		Info:          service.Info,
		ImgPath:       path,
		Price:         service.Price,
		DiscountPrice: service.DiscountPrice,
		OnSale:        true,
		Num:           service.Num,
		BossId:        uid,
		BossName:      boss.UserName,
		BossAvatar:    boss.Avatar,
	}

	// 将创建的 Product 插入数据库
	productDao := dao.NewProductDao(ctx)
	err = productDao.CreateProduct(product)
	if err != nil {
		code = e.Error
		utils.LogrusObj.Infoln("ProductService func Create:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 接下来需要将 files 包含的若干商品图片都上传到本地相应路径，并插入到 ProductImg 数据库里
	// 因为图片包含多张，所以采用并发处理
	// 创建 sync.WaitGroup 实例，用于等待所有图片上传完成。
	// wg.Add(len(files)) 将 WaitGroup 计数器设置为图片总数
	wg := new(sync.WaitGroup)
	wg.Add(len(files))
	for index, file := range files { // 遍历每张图片
		num := strconv.Itoa(index) // 使用 strconv.Itoa(index) 将索引转换为字符串，拼接在文件名中，确保文件名唯一

		// 不使用 productImgDao := dao.NewProductImgDao(ctx)
		// 这里使用了 productDao.DB，意味着 ProductImgDao 将复用 productDao 的数据库连接。
		// 这种方式的优势在于：
		// 1.共享事务环境（如果 productDao.DB 是事务实例，则 ProductImgDao 操作也会纳入该事务）。
		//   若商品创建成功但某张图片上传失败，使用同一事务的数据库连接可以保证所有操作自动回滚，避免出现“商品存在但图片缺失”的数据不一致问题
		// 2.避免额外的数据库连接开销，提升性能.
		//   当 files 包含大量图片时，每次新建连接会产生额外开销（如 TCP 握手、连接池竞争）。复用现有连接能提升性能，尤其在并发场景下。
		productImgDao := dao.NewProductImgDaoByDB(productDao.DB)

		// 打开当前图片文件并上传到本地静态目录
		// service.Name + num 作为文件名，即商品名+图片编号
		tmp, _ = file.Open()
		path, err = UploadProductToLocalStatic(tmp, uid, service.Name+num)
		if err != nil {
			code = e.ErrorProductImgUpload
			utils.LogrusObj.Infoln("ProductService func Create:", err)
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		// 构造 productImg 数据对象
		productImg := model.ProductImg{
			ProductId: product.ID,
			ImgPath:   path,
		}
		// 将这个 productImg 插入到数据库中
		err = productImgDao.CreateProductImg(&productImg)
		if err != nil {
			code = e.Error
			utils.LogrusObj.Infoln("ProductService func Create:", err)
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		// 成功后调用 wg.Done() 标记当前图片上传完成
		wg.Done()
	}
	// 等待所有图片上传完成后，继续执行
	wg.Wait()
	utils.LogrusObj.Infof("ProductService load %d product images OK.\n", len(files))
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildProduct(ctx, product), // 返回给前端 productVO 数据
	}
}

// List 获取商品列表
func (service *ProductService) List(ctx context.Context) serializer.Response {
	var products []*model.Product
	var err error

	code := e.Success
	// 如果未指定 `PageSize`，默认每页 15 条数据
	if service.PageSize == 0 {
		service.PageSize = 15
	}

	// condition 描述查询条件
	condition := make(map[string]interface{})
	if service.CategoryId != 0 {
		// 如果传入了 `CategoryId`，添加到查询条件中，展示的商品只会展示相应 CategoryId 的商品
		condition["category_id"] = service.CategoryId
	}

	productDao := dao.NewProductDao(ctx)
	// 查询满足条件的商品总数 total
	total, err := productDao.CountProductByCondition(condition)
	if err != nil {
		code = e.Error
		utils.LogrusObj.Infoln("ProductService CountProductByCondition:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 创建一个 WaitGroup 用于等待 Goroutine 完成
	// 与直接开启协程的区别是；
	// 使用 WaitGroup 保证当前协程等待 Goroutine 完成后再继续执行。
	// 如果不使用 WaitGroup，由于 Goroutine 异步执行，不阻塞主协程，主协程不会等待 products 数据获取完成，可能导致 products 数据未准备好时继续向下执行
	// 在获取数据后立即使用的场景下，比如这里获取 products 就要用 BuildProducts 来构建 productsVO，就要用 WaitGroup
	// TODO: go func 是否有作用？
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		productDao = dao.NewProductDaoByDB(productDao.DB)
		products, _ = productDao.ListProductByCondition(condition, service.BasePage)
		wg.Done()
	}()
	wg.Wait()
	return serializer.BuildListResponse(serializer.BuildProducts(ctx, products), uint(total))
}
