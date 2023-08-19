package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
//var usersLoginInfo = map[string]User{
//	"zhangleidouyin": {
//		Id:            1,
//		Name:          "zhanglei",
//		FollowCount:   10,
//		FollowerCount: 5,
//		IsFollow:      true,
//	},
//}

//var userIdSequence = int64(1)

type UserLoginResponse struct {
	model.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	model.Response
	User model.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	// 检验数据是否存在
	user := User{}
	if res := db.Where("token = ?", token).First(&user); res.Error == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{1, "User already exist"},
		})
	} else { // 不存在
		mutex.Lock()

		var userCount int64
		db.Model(&user).Count(&userCount)
		newUser := User{Id: userCount + 1, Name: username, Token: token}
		db.Create(&newUser)
		mutex.Unlock()

		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userCount + 1,
			Token:    token,
		})
	}

	/*if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := model.User{
			Id:   userIdSequence,
			Name: username,
		}
		usersLoginInfo[token] = newUser
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    username + password,
		})
	}*/
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	user := User{}
	if res := db.Where("token = ?", token).First(&user); res.Error == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在或密码错误"},
		})
	}

	/*if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}*/
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	user := User{}
	if res := db.Where("token = ?", token).First(&user); res.Error == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
	/*if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}*/
}
