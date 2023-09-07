package service

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
)

func InitFavoriteVideo(videos []model.Video, userId int64) []model.Video {
	favorite := model.FavoriteData{}
	// 获取用户点赞视频，并把IsFavorite改为true
	for i := range videos {
		favorite.IsFavorite = false
		db.DB.Where("video_id = ? and user_id = ?", videos[i].Id, userId).Find(&favorite)
		videos[i].IsFavorite = favorite.IsFavorite
	}
	return videos
}
