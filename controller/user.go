package controller

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
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

// var userIdSequence = int64(1)
var mutex sync.Mutex

type UserLoginResponse struct {
	model.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	model.Response
	User model.User `json:"user"`
}

// 后续修改
func generateToken(username string) (token string) {
	return username
}
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 检验数据是否存在
	user := model.User{}
	if res := db.DB.Where("id = ?", username).First(&user); res.Error == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else { // 不存在
		mutex.Lock()

		//var userCount int64
		//db.DB.Model(&user).Count(&userCount)
		//id数据库自增实现更好
		newUser := model.User{Username: username, Name: username, Password: password}
		db.DB.Create(&newUser)
		db.DB.Last(&newUser)
		mutex.Unlock()

		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 0},
			UserId:   newUser.Id,
			Token:    generateToken(newUser.Username),
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

	//token := username + password

	user := model.User{}
	if res := db.DB.Where("username = ?", username).Where("password=?", password).First(&user); res.Error == nil {
		token := generateToken(username)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "用户不存在或密码错误"},
		})
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
	user := model.User{}
	//为了测试暂时将token改为username了，后续引入redis再作讨论
	if res := db.DB.Where("username = ?", token).First(&user); res.Error == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
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
