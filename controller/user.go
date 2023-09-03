package controller

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
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
	user := model.User{}
	if res := db.DB.Where("username = ?", username).First(&user); res.Error == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else { // 不存在
		mutex.Lock()

		// encrypted : 已加密的密码
		encrypted, _ := GetPassword(password)
		//var userCount int64
		//db.DB.Model(&user).Count(&userCount)
		//id数据库自增实现更好

		newUser := model.User{Username: username, Name: username, Password: encrypted}
		db.DB.Create(&newUser)
		db.DB.Last(&newUser)
		mutex.Unlock()
		token, err := CreateToken(newUser.Id, newUser.Username)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 0, StatusMsg: "Registration successful"},
			UserId:   newUser.Id,
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

	//token := username + password

	user := model.User{}
	if res := db.DB.Where("username = ?", username).First(&user); res.Error != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
		return
	}
	token := generateToken(username)
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
	user := model.User{}
	claims, err := ParseToken(token)
	if err != nil {
		log.Fatal(err)
	}
	if res := db.DB.Where("id = ?", claims.UserId).First(&user); res.Error == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 0, StatusMsg: "Query success"},
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
