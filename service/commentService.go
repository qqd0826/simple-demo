package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/model"
	"strconv"
	"time"
)

func AddComment(videoId int64, user model.User, text string) model.Comment {
	comment := dao.AddComment(videoId, user, text)
	dao.VideoCommentCountInt(videoId)
	return comment
}
func DeleteComment(userId int64, commentId int) {
	// 检查评论用户ID和当前ID是否一致
	comment := dao.GetCommentById(commentId)
	if comment.UserId == userId {
		dao.DeleteComment(comment)
		// 更新视频的评论数
		dao.VideoCommentCountDec(comment.VideoId)
	}
}
func GetLastCommentList(videoId int64) []model.Comment {
	comments := dao.GetLastCommentList(videoId)
	// 日期格式化
	for i := range comments {
		data, _ := strconv.Atoi(comments[i].CreateDate)
		comments[i].CreateDate = time.Unix(int64(data), 0).Format("2006-01-02 15:04:05")
	}
	return comments
}
