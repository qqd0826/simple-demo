package dao

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"strconv"
	"time"
)

func AddComment(videoId int64, user model.User, text string) model.Comment {
	comment := model.Comment{
		VideoId:    videoId,
		User:       user,
		UserId:     user.Id,
		Content:    text,
		CreateDate: strconv.Itoa(int(time.Now().Unix())),
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
	return comments
}
