package service

import (
	"context"
	"gin-mail/dao"
	"gin-mail/model"
	"gin-mail/pkg/e"
	"gin-mail/serializer"
	"strconv"
)

type AddressService struct {
	Name    string `json:"name" form:"name"`
	Phone   string `json:"phone" form:"phone"`
	Address string `json:"address" form:"address"`
}

func (service *AddressService) Create(ctx context.Context, uid uint) serializer.Response {
	var address *model.Address
	code := e.Success

	addressDao := dao.NewAddressDao(ctx)

	address = &model.Address{
		UserID:  uid,
		Name:    service.Name,
		Phone:   service.Phone,
		Address: service.Address,
	}

	// 在数据库中插入 address 记录
	err := addressDao.CreateAddress(address)
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
	}
}

func (service *AddressService) Show(ctx context.Context, aid string) serializer.Response {
	addressId, _ := strconv.Atoi(aid)
	var address *model.Address
	code := e.Success

	addressDao := dao.NewAddressDao(ctx)

	// 在数据库中获取 address 记录
	address, err := addressDao.GetAddressById(uint(addressId))
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
		Data:   serializer.BuildAddress(address),
	}
}

func (service *AddressService) List(ctx context.Context, uid uint) serializer.Response {
	code := e.Success

	addressDao := dao.NewAddressDao(ctx)

	// 在数据库中获取用户的所有 address 记录
	addressList, err := addressDao.ListAddressByUserId(uid)
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
		Data:   serializer.BuildAddresses(addressList),
	}
}

func (service *AddressService) Update(ctx context.Context, uid uint, aid string) serializer.Response {
	addressId, _ := strconv.Atoi(aid)
	var address *model.Address
	code := e.Success

	addressDao := dao.NewAddressDao(ctx)

	address = &model.Address{
		Name:    service.Name,
		Phone:   service.Phone,
		Address: service.Address,
	}

	// 在数据库中更新 address 记录
	updated, err := addressDao.UpdateAddressByIdAndUserId(uint(addressId), uid, address)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if !updated { // 没有匹配的 address 记录，修改失败
		code = e.ErrorAddressNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *AddressService) Delete(ctx context.Context, uid uint, aid string) serializer.Response {
	addressId, _ := strconv.Atoi(aid)
	code := e.Success

	addressDao := dao.NewAddressDao(ctx)

	// 在数据库中删除 address 记录
	deleted, err := addressDao.DeleteAddressByIdAndUserId(uint(addressId), uid)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if !deleted { // 没有匹配的 address 记录，删除失败
		code = e.ErrorAddressNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
