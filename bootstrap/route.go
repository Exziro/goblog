package bootsrap

import (
	"goblog/app/models"
	"goblog/pkg/route"
	"goblog/routes"
	"time"

	"github.com/gorilla/mux"
)

func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)
	route.SetRoute(router)
	return router
}

//setupDB初始化数据库和 ORM
func SetupDB() {
	//建立数据库连接池
	db := models.ConnectDB()
	//命令行打出数据库请求信息
	sqlDB, _ := db.DB()

	// 设置最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}
