package db

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "gorm.io/driver/mysql"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func InitDB() {
	connect := "root:123456@tcp(localhost:3306)/godb?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", connect)
	if err != nil {
		log.Fatal(err)
	}

	//connect := "server=127.0.0.1;port=1433;database=dy;userid=sa;password=963013"
	//db, err := gorm.Open("mssql", connect)
	//if err != nil {
	//	log.Fatal(err)
	//}
	DB = db
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Video{})
	DB.AutoMigrate(&model.FavoriteData{})
	DB.AutoMigrate(&model.Comment{})
	// 这里可以继续添加其他需要迁移的表
	// 若不存在数据库表会自动生成
	DB.AutoMigrate(&model.Video{})
}
func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		log.Fatal(err)
	}
}
func Likes(userId int64, videoId int64) {
	val, _ := Redis.HGet(strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10)).Result()
	if len(val) == 0 {
		Redis.HSet(strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10), strconv.FormatInt(1, 10)+":"+strconv.FormatInt(time.Now().Unix(), 10))
		fmt.Println(strconv.FormatInt(1, 10) + ":" + strconv.FormatInt(time.Now().Unix(), 10))
	} else {
		//Redis.HSet(strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10), strconv.FormatInt(-1, 10)+":"+strconv.FormatInt(time.Now().Unix(), 10))
		Redis.HDel(strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10))
		fmt.Println(strconv.FormatInt(-1, 10) + ":" + strconv.FormatInt(time.Now().Unix(), 10))
	}

}
func CancelLikes(userId int64, videoId int64) {
	val, _ := Redis.HGet(strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10)).Result()
	if len(val) == 0 {
		Redis.HSet(strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10), strconv.FormatInt(0, 10)+":"+strconv.FormatInt(time.Now().Unix(), 10))
		fmt.Println(strconv.FormatInt(0, 10) + ":" + strconv.FormatInt(time.Now().Unix(), 10))
	} else {
		Redis.HDel(strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10))
		//Redis.HSet(strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10), strconv.FormatInt(-1, 10)+":"+strconv.FormatInt(time.Now().Unix(), 10))
		fmt.Println(strconv.FormatInt(-1, 10) + ":" + strconv.FormatInt(time.Now().Unix(), 10))
	}
}
func GetLikeVideoIds(userId int64) []int64 {
	videoIds := make([]int64, 0, 10)
	fields, err := Redis.HGetAll(strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(fields)
	for field, value := range fields {
		valueArr := strings.Split(value, ":")
		intValue, _ := strconv.ParseInt(valueArr[0], 10, 64)
		if intValue == 1 {
			intField, _ := strconv.ParseInt(field, 10, 64)
			videoIds = append(videoIds, intField)
		}
	}
	return videoIds
}
