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

type CommentListResponse struct {
	model.Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	model.Response
	Comment model.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	actionType := c.Query("action_type")

	// 检验用户是否存在
	user := model.User{}
	if res := db.DB.Where("username = ?", token).First(&user); res.Error == nil {
		// 评论
		if actionType == "1" {
			text := c.Query("comment_text")

			// 插入数据
			comment := model.Comment{User: user, UserId: user.Id, VideoId: int64(videoId), Content: text, CreateDate: strconv.Itoa(int(time.Now().Unix()))}
			db.DB.Create(&comment)
			db.DB.Model(&model.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + 1"))

			c.JSON(http.StatusOK, CommentActionResponse{Response: model.Response{StatusCode: 0},
				Comment: comment})
		} else if actionType == "2" { // 删除评论
			commentId, _ := strconv.Atoi(c.Query("comment_id"))
			comment := model.Comment{}
			db.DB.Where("id = ?", commentId).First(&comment)
			if comment.UserId == user.Id {
				db.DB.Delete(&comment)
				db.DB.Model(&model.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1"))
			}
		}

		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "用户未登录，请先登录"})
	}

	/*if user, exist := usersLoginInfo[token]; exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			c.JSON(http.StatusOK, CommentActionResponse{Response: model.Response{StatusCode: 0},
				Comment: model.Comment{
					Id:         1,
					User:       user,
					Content:    text,
					CreateDate: "05-01",
				}})
			return
		}
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}*/
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	//token := c.Query("token")
	videoId, _ := strconv.Atoi(c.Query("video_id"))

	comments := []model.Comment{}
	db.DB.Order("create_date desc").Where("video_id = ?", videoId).Find(&comments)

	// 日期格式化
	for i := range comments {
		data, _ := strconv.Atoi(comments[i].CreateDate)
		comments[i].CreateDate = time.Unix(int64(data), 0).Format("2006-01-02 15:04:05")
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    model.Response{StatusCode: 0},
		CommentList: comments,
	})
}
