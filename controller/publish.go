package controller

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type VideoListResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	user := util.GetUserByToken(token)
	if user.Id == 0 {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "请重新登录"})
		return
	}

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
	fileName := service.UploadVideo(user, title, file)
	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  fileName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	token := c.Query("token")
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	//查找token是否存在
	user := util.GetUserByToken(token)
	if user.Id == 0 {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist-publish"})
		return
	}

	videos := service.GetUserLastVideoList(userId)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
