package util

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/robfig/cron/v3"
	"strconv"
	"sync"
	"time"
)

var (
	once    sync.Once
	rwMutex sync.RWMutex
)

func SyncFavorite() {
	c := cron.New()
	_, err := c.AddFunc("* 0/1 * * * ? *", func() {
		Sync()
	})
	if err != nil {
		return
	}
	c.Start()

}
func getRedisLock(key string) bool {
	// 尝试在 Redis 中设置锁
	result, err := db.Redis.SetNX(key, 1, time.Second*10).Result()
	if err != nil {
		fmt.Println("Failed to set lock in Redis:", err)
		return false
	}

	// 如果设置成功，则获取锁成功
	if result {
		return true
	}

	// 如果设置失败，则获取锁失败
	return false
}

func releaseRedisLock(key string) {
	// 释放锁，将 Redis 中的键删除
	db.Redis.Del(key).Val()
}

func readAllFromRedis() []model.FavoriteData {
	// 读操作开始，获取读锁
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	keys, err := db.Redis.Keys("*").Result()
	if err != nil {
		fmt.Println("获取键失败:", err)
		return nil
	}
	var favoriteDatas []model.FavoriteData

	for _, key := range keys {
		fields, _ := db.Redis.HGetAll(key).Result()
		userId, _ := strconv.ParseInt(key, 10, 64)
		favoriteDatas = append(favoriteDatas, dao.TurnToFavorite(fields, userId)...)
	}
	db.Redis.FlushAll()
	return favoriteDatas
}

func writeToRedis(userId int64, videoId int64, val bool) {
	// 写操作开始，获取写锁
	rwMutex.Lock()
	defer rwMutex.Unlock()

	// 向 Redis 中写入数据
	if val {
		db.Likes(userId, videoId)
	} else {
		db.CancelLikes(userId, videoId)
	}
}
func Write(userId int64, videoId int64, val bool) {
	if getRedisLock("write") {
		writeToRedis(userId, videoId, val)
		releaseRedisLock("write")
	} else {
		//fmt.Println("Failed to get lock from Redis, retrying...")
		time.Sleep(time.Second * 5) // 等待一段时间后再重试获取锁
		Write(userId, videoId, val) // 重试获取锁并执行读写操作
	}
}
func Sync() {
	var favoriteDatas []model.FavoriteData
	if getRedisLock("write") {
		favoriteDatas = readAllFromRedis()
		releaseRedisLock("write")
	} else {
		//fmt.Println("Failed to get lock from Redis, retrying...")
		time.Sleep(time.Second * 5) // 等待一段时间后再重试获取锁
		Sync()                      // 重试获取锁并执行读写操作
	}
	if len(favoriteDatas) != 0 {
		dao.UpdateFavoriteBatch(favoriteDatas)
		fmt.Println("sync!")
	}
}
