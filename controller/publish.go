package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var baseUrl = "https://simple-douyin-go.oss-cn-hangzhou.aliyuncs.com/"

type VideoListResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	claims, err := ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "请重新登录"})
		return
	}
	fmt.Print(claims.UserId)
	/*if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}*/

	// 接收数据
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	title := c.PostForm("title")
	uploadFileName, err := db.UploadHandler(claims.UserId, title, data)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "上传视频失败，请稍后再试",
		})
		log.Fatal(err)
		return

	} else {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 0,
			StatusMsg:  title + " uploaded successfully",
		})
		//保存数据
		//filename := filepath.Base(data.Filename)
		//finalName := fmt.Sprintf("%d_%s", claims.UserId, filename)
		//
		//saveFile := filepath.Join("./public/", finalName)
		//if err := c.SaveUploadedFile(data, saveFile); err != nil {
		//	c.JSON(http.StatusOK, model.Response{
		//		StatusCode: 1,
		//		StatusMsg:  err.Error(),
		//	})
		//	return
		//}
		var user model.User
		db.DB.Where("id=?", claims.UserId).First(&user)
		// 保存到数据库
		newVideo := model.Video{
			Author:     user,
			PlayUrl:    baseUrl + uploadFileName,
			CoverUrl:   baseUrl + uploadFileName,
			UpLoadTime: time.Now().Unix(),
		}
		db.DB.Create(&newVideo)
	}
	return
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: getVideo(),
	})
}
