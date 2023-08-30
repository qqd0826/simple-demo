package controller

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	actionType := c.Query("action_type")

	// 检验用户是否存在
	user := model.User{}
	if res := db.DB.Where("username = ?", token).First(&user); res.Error == nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})

		var video model.Video
		db.DB.Where("id = ?", int64(videoId)).First(&video)

		favoriteData := model.FavoriteData{UserId: user.Id, VideoId: int64(videoId)}
		if db.DB.Where("user_id = ? and video_id = ?", user.Id, video.Id).Find(&favoriteData).RecordNotFound() {
			db.DB.Create(&favoriteData)
		}

		// 点赞
		if actionType == "1" {
			// 更新两张表
			db.DB.Model(&video).Update("favorite_count", gorm.Expr("favorite_count + 1"))
			db.DB.Model(&favoriteData).Where("user_id = ? and video_id = ?", user.Id, video.Id).Updates(model.FavoriteData{IsFavorite: true, Time: time.Now().Unix()})
		} else if actionType == "2" { // 取消点赞
			db.DB.Model(&video).Update("favorite_count", gorm.Expr("favorite_count - 1"))
			db.DB.Model(&favoriteData).Where("user_id = ? and video_id = ?", user.Id, video.Id).Updates(map[string]interface{}{"user_id": user.Id, "video_id": video.Id, "IsFavorite": false, "Time": time.Now().Unix()})
		}
	} else { // 不存在
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "用户未登录，请先登录"})
	}

	//if _, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	//} else {
	//	c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	token := c.Query("token")

	//查找token是否存在
	user := model.User{}
	if db.DB.Where("username = ?", token).First(&user).RecordNotFound() {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	// 获取当前用户点赞信息
	videos := make([]model.Video, 0)
	favoriteData := make([]model.FavoriteData, 0)
	db.DB.Where("is_favorite = ? and user_id = ?", true, user.Id).Find(&favoriteData)

	// 获取点赞视频的ID
	videoIds := make([]int64, len(favoriteData))
	for i := range favoriteData {
		videoIds[i] = favoriteData[i].VideoId
	}

	// 查找对应视频
	db.DB.Find(&videos, videoIds)

	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
