package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
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
	user := model.User{}
	if res := db.DB.Where("id = ?", claims.UserId).First(&user); res.Error != nil {
		fmt.Println(res.Error)
		return
	}
	/*if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}*/

	// 接收数据
	file, err := c.FormFile("data")
	title := c.PostForm("title")
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	fileName, err := db.UploadHandler(user.Id, title, file)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	//获得oss存储中的url，以及视频抽帧功能获得第一毫秒的帧作为封面。
	playurl := fmt.Sprintf("%s%s", "https://qqd-simple-demo.oss-cn-beijing.aliyuncs.com/", fileName)
	coverurl := fmt.Sprintf("%s%s", playurl, "?x-oss-process=video/snapshot,t_0001,f_jpg,w_800,h_600,m_fast")

	// 保存到数据库
	newVideo := model.Video{
		AuthorId:   user.Id,
		Author:     user,
		PlayUrl:    playurl,
		CoverUrl:   coverurl,
		UpLoadTime: time.Now().Unix(),
		Title:      title,
	}

	db.DB.Create(&newVideo)
	// 更新用户表中的视频数量
	db.DB.Model(&user).Update("work_count", gorm.Expr("work_count + 1"))

	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  fileName + " uploaded successfully",
	})

	//// 保存数据
	//filename := filepath.Base(data.Filename)
	//finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	//saveFile := filepath.Join("./public/", finalName)
	//if err := c.SaveUploadedFile(data, saveFile); err != nil {
	//	c.JSON(http.StatusOK, model.Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}
	//// 保存到数据库
	//newVideo := model.Video{
	//	Author:     user,
	//	PlayUrl:    saveFile,
	//	CoverUrl:   saveFile,
	//	UpLoadTime: time.Now().Unix(),
	//	Title:      title,
	//}
	//db.DB.Create(&newVideo)
	//
	//c.JSON(http.StatusOK, model.Response{
	//	StatusCode: 0,
	//	StatusMsg:  finalName + " uploaded successfully",
	//})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")

	//查找token是否存在
	user := model.User{}
	res := db.DB.Where("username = ?", token).First(&user)
	if res.Error != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	videos := []model.Video{}

	db.DB.Order("up_load_time desc").Where("author_id = ?", user_id).Find(&videos)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
