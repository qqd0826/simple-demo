package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/model"
)

func AddUserThenGet(user model.User) model.User {
	dao.AddUser(user)
	dao.GetUser(user)
	return user
}
