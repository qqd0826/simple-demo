package db

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "gorm.io/driver/mysql"
	"log"
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
	if err != nil {
		log.Fatal(err)
	}
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
