package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]model.User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	model.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	model.Response
	User model.User `json:"user"`
}

// GetPassword 给密码加密
func GetPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

// CheckPassword 用于比对密码和哈希值是否匹配,如果匹配，则返回 true，否则返回 false。
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	return err == nil
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 检验数据是否存在
	user := dao.GetUserByUsername(username)
	if user.Id != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else { // 不存在
		encrypted, _ := GetPassword(password)
		newUser := model.User{Username: username, Name: username, Password: encrypted}
		newUser = service.AddUserThenGet(newUser)
		token, err := util.CreateToken(newUser.Id, newUser.Username)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 0, StatusMsg: "Registration successful"},
			UserId:   newUser.Id,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user := dao.GetUserByUsername(username)
	if user.Id == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
		return
	}
	token, err := util.CreateToken(user.Id, username)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(token)
	if CheckPassword(user.Password, password) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 0, StatusMsg: "登录成功"},
			UserId:   user.Id,
			Token:    token,
		})
		return
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "密码错误"},
		})
		return
	}

	//if user, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: model.Response{StatusCode: 0},
	//		UserId:   user.Id,
	//		Token:    token,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	user := util.GetUserByToken(token)
	if user.Id != 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 0, StatusMsg: "Query success"},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist-user"},
		})
	}
}
