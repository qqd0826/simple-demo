package db

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "gorm.io/driver/mysql"
	"log"
)

var DB *gorm.DB

func InitDB() {
	/*connect := "root:123456@tcp(localhost:3306)/godb?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", connect)
	if err != nil {
		log.Fatal(err)
	}*/

	connect := "server=127.0.0.1;port=1433;database=dy;user id=sa;password=963013"
	db, err := gorm.Open("mssql", connect)
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

}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		log.Fatal(err)
	}
}
