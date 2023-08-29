package controller

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")

	// 检验用户是否存在
	user := model.User{}
	if res := db.DB.Where("username = ?", token).First(&user); res.Error == nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})

		var video model.Video
		db.DB.Where("id = ?", videoId).First(&video)

		if actionType == "1" {
			db.DB.Model(&video).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
			db.DB.Model(&video).Update("is_favorite", true)
		} else if actionType == "2" {
			db.DB.Model(&video).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
			db.DB.Model(&video).Update("is_favorite", false)
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
	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
