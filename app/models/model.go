package models

import (
	"goblog/pkg/logger"
	"goblog/pkg/types"
	"time"

	//命令行查看调试语句
	gormlogger "gorm.io/gorm/logger"
	// GORM 的 MSYQL 数据库驱动导入
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//DB gorm.DB对象
var DB *gorm.DB

// BaseModel 模型基类
type BaseModel struct {
	ID uint64

	CreatedAt time.Time `gorm:"column:created_at;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;index"`
}

// GetStringID 获取 ID 的字符串格式
func (a BaseModel) GetStringID() string {
	return types.Uint64ToString(a.ID)
}

//connectDB初始化数据模型
func ConnectDB() *gorm.DB {
	var err error
	config := mysql.New(mysql.Config{
		DSN: "root:qiyixi19961016@tcp(127.0.0.1:3306)/goblog?charset=utf8&parseTime=True&loc=Local",
	})
	//准备数据链接池
	DB, err = gorm.Open(config, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	logger.LogError(err)

	return DB
}
