package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type VideoListResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")

	// 用户是否存在，目前token就是用户名
	user := model.User{}
	res := db.DB.Where("username = ?", token).First(&user)
	if res.Error != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
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
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以北京为例，填写为https://oss-cn-beijing.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New(
		"https://oss-cn-beijing.aliyuncs.com",
		"LTAI5tDPiHQY8DXVUv25mhQH",
		"ej9I1AAHc6OM7DP4W70W9l18074lrH",
	)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket("qqd-simple-demo")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	//保存到oss
	file, _ := data.Open()
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d-%v-%s", user.Id, time.Now().Unix(), filename)
	err = bucket.PutObject(finalName, file)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	//获得oss存储中的url，以及视频抽帧功能获得第一毫秒的帧作为封面。
	playurl := fmt.Sprintf("%s%s", "https://qqd-simple-demo.oss-cn-beijing.aliyuncs.com/", finalName)
	coverurl := fmt.Sprintf("%s%s", playurl, "?x-oss-process=video/snapshot,t_0001,f_jpg,w_800,h_600,m_fast")

	// 保存到数据库
	newVideo := model.Video{
		Author:     user,
		PlayUrl:    playurl,
		CoverUrl:   coverurl,
		UpLoadTime: time.Now().Unix(),
		Title:      title,
	}
	db.DB.Create(&newVideo)

	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
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
	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: getVideo(),
	})
}
