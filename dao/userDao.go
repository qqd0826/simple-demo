package dao

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/jinzhu/gorm"
)

func GetUserById(id int64) model.User {
	var user model.User
	db.DB.Where("id=?", id).First(&user)
	return user
}
func GetUserByUsername(username string) model.User {
	var user model.User
	db.DB.Where("username = ?", username).First(&user)
	return user
}
func AddUser(newUser model.User) {
	db.DB.Create(&newUser)
}
func GetUser(user model.User) model.User {
	db.DB.Last(&user)
	return user
}
func UserWorkCountInc(user model.User) {
	db.DB.Model(&user).Update("work_count", gorm.Expr("work_count + 1"))
}
