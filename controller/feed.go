package controller

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	token := c.Query("token")

	//如果未登录，直接返回未点赞的video列表
	user := model.User{}
	if db.DB.Where("username = ?", token).First(&user).RecordNotFound() {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  model.Response{StatusCode: 0},
			VideoList: getVideo(),
			NextTime:  time.Now().Unix(),
		})
		return
	}

	//feed流的video
	feedVideo := getVideo()

	favorite := model.FavoriteData{}
	// 获取用户点赞视频，并把IsFavorite改为true
	for i := range feedVideo {
		favorite.IsFavorite = false
		db.DB.Where("video_id = ? and user_id = ?", feedVideo[i].Id, user.Id).Find(&favorite)
		feedVideo[i].IsFavorite = favorite.IsFavorite
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: feedVideo,
		NextTime:  time.Now().Unix(),
	})
	return
}

func getVideo() (videos []model.Video) {
	//按投稿时间倒序，限制30个
	db.DB.Limit(30).Order("up_load_time desc").Find(&videos)
	return
}
