package dao

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
)

func GetFavorite(userId int64, videoId int64) model.FavoriteData {
	favoriteData := model.FavoriteData{}
	db.DB.Where("user_id = ? and video_id = ?", userId, videoId).Find(&favoriteData)
	return favoriteData
}
