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
func GetFavoritVideoList(userId int64) []model.Video {
	feedVideo := GetLastVideoList()
	return InitFavoriteVideo(feedVideo, userId)
}
