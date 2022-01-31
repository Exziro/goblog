package main

import (

	//"go/types"
	//"goblog/pkg/logger"
	"embed"
	"goblog/app/http/middlewares"
	bootsrap "goblog/bootstrap"
	"goblog/config"

	c "goblog/pkg/config"
	"goblog/pkg/logger"

	//"goblog/pkg/model"

	//"errors"

	//"html/template"
	"net/http"

	//"time"
	//"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//go:embed resources/views/articles/*
//go:embed resources/views/auth/*
//go:embed resources/views/categories/*
//go:embed resources/views/layouts/*
var tplFS embed.FS

//go:embed public/*
var staticFS embed.FS

var router *mux.Router

//博文存储部分函数
func init() {
	// 初始化配置信息
	config.Initialize()
}
func main() {

	//初始化SQL
	bootsrap.SetupDB()

	bootsrap.SetupTemplate(tplFS)
	//初始化路由
	router = bootsrap.SetupRoute(staticFS)
	err := http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)

}
