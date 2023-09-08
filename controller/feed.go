package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/util"
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
	fmt.Println(token)
	//如果未登录，直接返回未点赞的video列表
	if len(token) == 0 || util.GetUserByToken(token).Id == 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  model.Response{StatusCode: 0},
			VideoList: getVideo(),
			NextTime:  time.Now().Unix(),
		})
		return

	} else {
		user := util.GetUserByToken(token)
		feedVideo := service.GetFeedVideoList(user.Id)

		c.JSON(http.StatusOK, FeedResponse{
			Response:  model.Response{StatusCode: 0},
			VideoList: feedVideo,
			NextTime:  time.Now().Unix(),
		})
		return
	}
}

func getVideo() (videos []model.Video) {
	//按投稿时间倒序，限制30个
	service.GetLastVideoList()
	return
}
