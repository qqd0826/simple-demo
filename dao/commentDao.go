package dao

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"time"
)

func AddComment(videoId int64, user model.User, text string) model.Comment {
	comment := model.Comment{
		VideoId:    int64(videoId),
		User:       user,
		UserId:     user.Id,
		Content:    text,
		CreateDate: time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"),
	}
	db.DB.Create(&comment)
	return comment
}
func GetCommentById(id int) model.Comment {
	comment := model.Comment{}
	db.DB.Where("id = ?", id).First(&comment)
	return comment
}
func DeleteComment(comment model.Comment) {
	db.DB.Delete(&comment)
}
func GetLastCommentList(videoId int64) []model.Comment {
	var comments []model.Comment
	db.DB.Order("create_date desc").Where("video_id = ?", videoId).Find(&comments)
	for i := range comments {
		user := GetUserById(comments[i].UserId)
		user.Password = ""
		comments[i].User = user
	}
	return comments
}
