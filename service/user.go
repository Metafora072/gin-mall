package service

import (
	"context"
	"gin-mail/conf"
	"gin-mail/dao"
	"gin-mail/model"
	"gin-mail/pkg/e"
	"gin-mail/pkg/utils"
	"gin-mail/serializer"
	"gopkg.in/mail.v2"
	"mime/multipart"
	"strings"
)

// UserService 是用户普通请求绑定的结构体
type UserService struct {
	NickName string `json:"nick_name" form:"nick_name"`
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"` // 前端验证
}

// SendEmailService 是用户发送邮箱请求绑定的结构体
// OperationType: 操作类型：1.绑定邮箱；2.解绑邮箱；3.改密码
type SendEmailService struct {
	Email         string `json:"email" form:"email"`
	Password      string `json:"password" form:"password"`
	OperationType uint   `json:"operation_type" form:"operation_type"`
}

// Register 处理用户注册逻辑
func (service *UserService) Register(ctx context.Context) serializer.Response {
	var user model.User
	code := e.Success

	// 加密 Key 密钥长度要求是 16
	if service.Key == "" || len(service.Key) != 16 {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "密钥长度不足",
		}
	}

	// 对称加密
	utils.Encrypt.SetKey(service.Key)

	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if exist {
		code = e.ErrorExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	user = model.User{
		UserName: service.UserName,
		NickName: service.NickName,
		Status:   model.Active,
		Avatar:   "avatar.jpg",
		Money:    utils.Encrypt.AesEncoding("10000"), // 初始金额 10000, 对其进行加密
	}
	// 密码加密
	if err = user.SetPassword(service.Password); err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 创建用户
	err = userDao.CreateUser(&user)
	if err != nil {
		code = e.Error
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// Login 处理用户登录逻辑
func (service *UserService) Login(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	// 检验相应用户记录是否在数据库中存在
	if !exist || err != nil {
		code = e.ErrorExistUserNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户不存在, 请先注册",
		}
	}
	// 校验登录密码
	if user.CheckPassword(service.Password) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密码错误, 请重新登录",
		}
	}
	// 签发 token
	token, err := utils.GenerateToken(user.ID, user.UserName, 0)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密码错误, 请重新登录",
		}
	}
	return serializer.Response{
		Status: code,
		Data: serializer.TokenData{
			User:  serializer.BuildUser(user),
			Token: token,
		},
		Msg: e.GetMsg(code),
	}
}

// Update 处理用户修改信息逻辑(目前仅支持修改昵称 NickName)
func (service *UserService) Update(ctx context.Context, uid uint) serializer.Response {
	var user *model.User
	var err error
	code := e.Success
	// 在数据库中找到相应用户的记录
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(uid)
	// 修改昵称
	if service.NickName != "" {
		user.NickName = service.NickName
	}
	err = userDao.UpdateUserById(uid, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// Post 处理用户上传头像的逻辑
func (service *UserService) Post(ctx context.Context, uid uint, file multipart.File, fileSize int64) serializer.Response {
	code := e.Success
	var user *model.User
	var err error
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(uid)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	// 保存图片到本地
	path, err := UploadAvatarToLocalStatic(file, uid, user.UserName)
	if err != nil {
		code = e.ErrorUploadFail
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	user.Avatar = path
	err = userDao.UpdateUserById(uid, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}

}

// Send 处理发送邮箱逻辑
func (service *SendEmailService) Send(ctx context.Context, uid uint) serializer.Response {
	code := e.Success
	var address string
	var notice *model.Notice

	// 生成包含用户ID、操作类型、邮箱和密码的加密令牌
	// 该令牌将作为验证链接参数，后续需要解密验证
	// 注意：此处直接将密码传入可能存在安全隐患，建议使用加密传输
	emailToken, err := utils.GenerateEmailToken(uid, service.OperationType, service.Email, service.Password)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 根据操作类型从数据库获取对应的邮件模板
	// notice表结构示例：
	// | id | 其他字段 | 模板内容                 |
	// |----|---------|-------------------------|
	// | 1  | ...     | 您正在绑定邮箱Email       |
	// | 2  | ...     | 您正在解绑邮箱Email       |
	// | 3  | ...     | 您正在修改密码Email       |
	noticeDao := dao.NewNoticeDao(ctx)
	notice, err = noticeDao.GetNoticeById(service.OperationType)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 拼接完整的验证地址
	// conf.ValidEmail = http://localhost:8080/
	address = conf.ValidEmail + emailToken // 发送方

	// 替换模板中的 Email 占位符为实际验证链接:
	// 您正在绑定邮箱 Email 中的 Email 占位符被替换为 address
	mailStr := notice.Text
	mailText := strings.Replace(mailStr, "Email", address, -1)

	m := mail.NewMessage()
	// 设置邮件头信息
	m.SetHeader("From", conf.SmtpEmail) // 发件人（从config.ini配置中读取的邮箱）
	m.SetHeader("To", service.Email)    // 收件人（用户提交的邮箱）
	m.SetHeader("Subject", "Metafora")  // 邮件主题
	m.SetBody("text/html", mailText)    // HTML格式正文

	// 建立SMTP连接并发送
	// 创建SMTP拨号器（TLS加密）
	d := mail.NewDialer(
		conf.SmtpHost,  // SMTP服务器地址
		465,            // SSL端口
		conf.SmtpEmail, // SMTP账号
		conf.SmtpPass,  // SMTP密码/授权码
	)
	d.StartTLSPolicy = mail.MandatoryStartTLS // 强制TLS加密

	// DialAndSend 打开与 SMTP 服务器的连接，发送给定的电子邮件并关闭连接。
	if err = d.DialAndSend(m); err != nil {
		code = e.ErrorSendEmail
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
