package controller

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType := c.Query("action_type")

	// 检验用户是否存在
	user := util.GetUserByToken(token)
	if user.Id == 0 {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "用户登录信息失效，请重新登录"})
	}

	// 评论
	if actionType == "1" {
		text := c.Query("comment_text")

		// 插入数据并更新视频的评论数
		comment := service.AddComment(videoId, user, text)

		c.JSON(http.StatusOK, CommentActionResponse{Response: model.Response{StatusCode: 0},
			Comment: comment,
		})
		return
	} else if actionType == "2" {
		//删除评论
		commentId, _ := strconv.Atoi(c.Query("comment_id"))
		service.DeleteComment(user.Id, commentId)
		c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: "评论删除成功"})
		return
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
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)

	comments := service.GetLastCommentList(videoId)

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    model.Response{StatusCode: 0},
		CommentList: comments,
	})
}
