package controller

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType := c.Query("action_type")

	// 检验用户是否存在
	user := util.GetUserByToken(token)
	if user.Id != 0 {
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
		//video := service.GetVideoById(videoId)
		////favoriteData := model.FavoriteData{UserId: user.Id, VideoId: videoId}
		////if db.DB.Where("user_id = ? and video_id = ?", user.Id, video.Id).Find(&favoriteData).RecordNotFound() {
		////	db.DB.Create(&favoriteData)
		////}

		// 点赞
		if actionType == "1" {
			// 更新两张表
			//db.DB.Model(&video).Update("favorite_count", gorm.Expr("favorite_count + 1"))
			//db.DB.Model(&favoriteData).Where("user_id = ? and video_id = ?", user.Id, video.Id).Updates(model.FavoriteData{IsFavorite: true, Time: time.Now().Unix()})
			//db.DB.Model(&user).Update("favorite_count", gorm.Expr("favorite_count + 1"))
			util.Write(user.Id, videoId, true)
		} else if actionType == "2" { // 取消点赞
			//db.DB.Model(&video).Update("favorite_count", gorm.Expr("favorite_count - 1"))
			//db.DB.Model(&favoriteData).Where("user_id = ? and video_id = ?", user.Id, video.Id).Updates(map[string]interface{}{"user_id": user.Id, "video_id": video.Id, "IsFavorite": false, "Time": time.Now().Unix()})
			//db.DB.Model(&user).Update("favorite_count", gorm.Expr("favorite_count - 1"))
			util.Write(user.Id, videoId, false)
		}
	} else { // 不存在
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "用户登录信息失效，请重新登录"})
	}
}

func FavoriteList(c *gin.Context) {
	//查找token是否存在
	token := c.Query("token")
	user := util.GetUserByToken(token)
	if user.Id == 0 {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "请重新登录"})
		return
	}

	// 获取当前用户点赞信息
	videos := service.GetFavoriteVideoList(user.Id)

	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
