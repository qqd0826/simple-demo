package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
)

func InitFavoriteVideo(videos []model.Video, userId int64) []model.Video {
	// 获取用户点赞视频，并把IsFavorite改为true
	idMap := make(map[int64]int)
	for i := range videos {
		idMap[videos[i].Id] = i
	}
	videoIds := getVideoIdsByVideos(videos)
	var favorite []model.FavoriteData
	db.DB.Where("video_id in (?) and user_id = ?", videoIds, userId).Find(&favorite)
	for i := range favorite {
		value, ok := idMap[favorite[i].VideoId]
		if ok {
			videos[value].IsFavorite = favorite[i].IsFavorite
		}
	}
	tempData := dao.GetLikeInfo(userId)
	for i := range tempData {
		for j := range videos {
			if tempData[i].VideoId == videos[j].Id {
				videos[j].IsFavorite = tempData[i].IsFavorite
				if tempData[i].IsFavorite {
					videos[j].FavoriteCount++
				}
			}
		}
	}
	return videos
}
