package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/util"
	"mime/multipart"
	"os"
	"time"
)

func AddVideo(user model.User, video model.Video) {
	dao.AddVideo(video)
	dao.UserWorkCountInc(user)
}
func UploadVideo(user model.User, title string, file *multipart.FileHeader) (fileName string) {
	fileName, err := db.UploadHandler(user.Id, title, file)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	//获得oss存储中的url，以及视频抽帧功能获得第一毫秒的帧作为封面。
	playurl := fmt.Sprintf("%s%s", util.BaseUrl, fileName)
	coverurl := fmt.Sprintf("%s%s", playurl, "?x-oss-process=video/snapshot,t_0001,f_jpg,w_800,h_600,m_fast")
	fmt.Println(playurl)
	fmt.Println(coverurl)
	// 保存到数据库
	newVideo := model.Video{
		AuthorId:   user.Id,
		Author:     user,
		PlayUrl:    playurl,
		CoverUrl:   coverurl,
		UpLoadTime: time.Now().Unix(),
		Title:      title,
	}
	AddVideo(user, newVideo)
	return fileName
}
func GetUserLastVideoList(userId int64) []model.Video {
	return dao.GetUserLastVideoList(userId)
}
func GetVideoById(id int64) model.Video {
	return dao.GetVideoById(id)

}
func GetLastVideoList() []model.Video {
	return dao.GetLastVideoList()
}
func GetFeedVideoList(userId int64) []model.Video {
	feedVideo := GetLastVideoList()
	return InitFavoriteVideo(feedVideo, userId)
}
func GetFavoriteVideoList(userId int64) []model.Video {

	favoriteData := dao.GetUserFavoriteData(userId)
	videoIds := getVideoIdsByFavorite(favoriteData)
	videos := dao.GetVideosByIdList(videoIds)
	fmt.Println(videoIds)
	favoriteDataFromRedis := dao.GetLikeInfo(userId)
	fmt.Print(len(favoriteDataFromRedis), favoriteDataFromRedis)
	allVideo := make([]model.Video, len(videos))
	duplicateFlag := make([]bool, len(favoriteDataFromRedis))
	for i := range videos {
		for j := range favoriteDataFromRedis {
			fmt.Println(j, favoriteDataFromRedis[j].VideoId)
			if videos[i].Id == favoriteDataFromRedis[j].VideoId {
				if favoriteDataFromRedis[j].IsFavorite {
					videos[i].FavoriteCount++
				} else {
					videos[i].FavoriteCount--
				}
				duplicateFlag[j] = true
			}
		}
	}
	allVideo = append(allVideo, videos...)
	for i := range favoriteDataFromRedis {
		if duplicateFlag[i] == false {
			video := dao.GetVideoById(favoriteDataFromRedis[i].VideoId)
			if favoriteDataFromRedis[i].IsFavorite {
				video.FavoriteCount++
			} else {
				video.FavoriteCount--
			}
			allVideo = append(allVideo, video)
		}

	}
	return allVideo
}
func getVideoIdsByVideos(videos []model.Video) []int64 {
	videoIds := make([]int64, len(videos))
	for i := range videos {
		videoIds[i] = videos[i].Id
	}
	return videoIds
}
func getVideoIdsByFavorite(favoriteData []model.FavoriteData) []int64 {
	videoIds := make([]int64, len(favoriteData))
	for i := range favoriteData {
		videoIds[i] = favoriteData[i].VideoId
	}
	return videoIds
}
func mergeAndDeduplicate(ids1 []int64, ids2 []int64) {
	result := make([]int64, len(ids1))
	result = append(result, ids1...)
	for j := range ids2 {
		flag := false
		for i := range result {
			if result[i] == ids2[j] {
				flag = true
				break
			}
		}
		if flag == false {
			result = append(result, ids2[j])
		}
	}
}
