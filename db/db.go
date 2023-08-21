package db

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "gorm.io/driver/mysql"
	"log"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func InitDB() {
	connect := "root:root@tcp(localhost:3306)/gorm?charset=utf8&parseTime=True&loc=Local"
	connect := "root:123456@tcp(localhost:3306)/godb?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", connect)
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	DB.AutoMigrate(&model.User{})
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
