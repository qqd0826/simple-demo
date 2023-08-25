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
	c.JSON(http.StatusOK, FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: getVideo(),
		NextTime:  time.Now().Unix(),
	})
}

func getVideo() (videos []model.Video) {
	//按投稿时间倒序，限制30个
	db.DB.Limit(30).Order("up_load_time desc").Find(&videos)
	return
}
