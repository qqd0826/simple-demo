package dao

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/jinzhu/gorm"
)

func AddVideo(newVideo model.Video) {
	db.DB.Create(&newVideo)
}
func GetUserLastVideoList(userId int64) []model.Video {
	videos := []model.Video{}
	db.DB.Order("up_load_time desc").Where("author_id = ?", userId).Find(&videos)
	return videos
}
func GetVideoById(id int64) model.Video {
	var video model.Video
	db.DB.Where("id = ?", id).First(&video)
	return video
}
func VideoCommentCountInt(videoId int64) {
	db.DB.Model(&model.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + 1"))
}
func VideoCommentCountDec(videoId int64) {
	db.DB.Model(&model.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1"))
}
func GetLastVideoList() []model.Video {
	var videos []model.Video
	db.DB.Limit(30).Order("up_load_time desc").Find(&videos)
	return videos
}
