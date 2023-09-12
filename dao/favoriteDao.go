package dao

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"strings"
)

var (
	flag = [2]bool{false, true}
)

func GetFavorite(userId int64, videoId int64) model.FavoriteData {
	favoriteData := model.FavoriteData{}
	db.DB.Where("user_id = ? and video_id = ?", userId, videoId).Find(&favoriteData)
	return favoriteData
}
func GetUserFavoriteData(userId int64) []model.FavoriteData {
	favoriteData := make([]model.FavoriteData, 0)
	db.DB.Where("is_favorite = ? and user_id = ?", true, userId).Find(&favoriteData)
	return favoriteData
}
func GetLikeInfo(userId int64) []model.FavoriteData {
	fields, err := db.Redis.HGetAll(strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		log.Println(err)
	}
	likeInfo := TurnToFavorite(fields, userId)
	return likeInfo
}
func TurnToFavorite(fields map[string]string, userId int64) []model.FavoriteData {
	favoriteDatas := make([]model.FavoriteData, 0)

	for field, value := range fields {
		valueArr := strings.Split(value, ":")
		intValue, _ := strconv.ParseInt(valueArr[0], 10, 64)
		if intValue != -1 {
			favoriteData := model.FavoriteData{}
			favoriteData.UserId = userId
			favoriteData.VideoId, _ = strconv.ParseInt(field, 10, 64)
			favoriteData.IsFavorite = flag[intValue]
			favoriteData.Time, _ = strconv.ParseInt(valueArr[1], 10, 64)
			favoriteDatas = append(favoriteDatas, favoriteData)
		}

	}
	fmt.Print(favoriteDatas)
	return favoriteDatas
}
func UpdateFavoriteBatch(favoriteDatas []model.FavoriteData) {
	for i := range favoriteDatas {
		temp := []model.FavoriteData{}
		item := favoriteDatas[i]
		db.DB.Where("user_id = ? and video_id = ?", item.UserId, item.VideoId).Find(&temp)
		if len(temp) == 0 {
			db.DB.Create(&item)
		} else {
			db.DB.Model(&model.FavoriteData{}).Where("user_id = ? and video_id = ?", item.UserId, item.VideoId).Updates(model.FavoriteData{IsFavorite: item.IsFavorite, Time: item.Time})
		}
		video := GetVideoById(item.VideoId)
		user := GetUserById(item.UserId)
		if item.IsFavorite {
			db.DB.Model(&video).Update("favorite_count", gorm.Expr("favorite_count + 1"))
			db.DB.Model(&user).Update("favorite_count", gorm.Expr("favorite_count + 1"))
		} else {
			db.DB.Model(&video).Update("favorite_count", gorm.Expr("favorite_count - 1"))
			db.DB.Model(&user).Update("favorite_count", gorm.Expr("favorite_count - 1"))
		}

	}
}

//func GetVideoLikeInfo(fields map[string]string, userId int64) []model.FavoriteData {
//	favoriteDatas := make([]model.FavoriteData, 0)
//	for field, value := range fields {
//		valueArr := strings.Split(value, ":")
//		intValue, _ := strconv.ParseInt(valueArr[0], 10, 64)
//		if intValue != -1 {
//			favoriteData := model.FavoriteData{}
//			favoriteData.UserId = userId
//			favoriteData.VideoId, _ = strconv.ParseInt(field, 10, 64)
//			favoriteData.IsFavorite = flag[intValue]
//			favoriteData.Time, _ = strconv.ParseInt(valueArr[1], 10, 64)
//			favoriteDatas = append(favoriteDatas, favoriteData)
//		}
//
//	}
//	return favoriteDatas
//}
