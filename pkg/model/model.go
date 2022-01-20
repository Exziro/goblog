package model

import (
	"goblog/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//DB gorm.DB对象
var DB *gorm.DB

//connectDB初始化数据模型
func ConnectDB() *gorm.DB {
	var err error
	config := mysql.New(mysql.Config{
		DSN: "root:qiyixi19961016@tcp(127.0.0.1:3306)/goblog?charset=utf8&parseTime=True&loc=Local",
	})
	//准备数据链接池
	DB, err := gorm.Open(config, &gorm.Config{})
	logger.LogError(err)

	return DB
}
