package models

import (
	"fmt"
	"goblog/pkg/config"
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

	// 初始化 MySQL 连接信息
	gormConfig := mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
			config.GetString("database.mysql.username"),
			config.GetString("database.mysql.password"),
			config.GetString("database.mysql.host"),
			config.GetString("database.mysql.port"),
			config.GetString("database.mysql.database"),
			config.GetString("database.mysql.charset")),
	})

	// gormConfig := mysql.New(mysql.Config{
	// 	DSN: dsn,
	// })

	var level gormlogger.LogLevel
	if config.GetBool("app.debug") {
		// 读取不到数据也会显示
		level = gormlogger.Warn
	} else {
		// 只有错误才会显示
		level = gormlogger.Error
	}

	// 准备数据库连接池
	DB, err = gorm.Open(gormConfig, &gorm.Config{
		Logger: gormlogger.Default.LogMode(level),
	})

	logger.LogError(err)

	return DB
}
