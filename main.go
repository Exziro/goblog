package main

import (
	"database/sql"
	//"go/types"
	//"goblog/pkg/logger"
	"goblog/app/http/middlewares"
	bootsrap "goblog/bootstrap"
	"goblog/pkg/logger"

	//"goblog/pkg/model"

	//"errors"

	//"html/template"
	"net/http"
	"net/url"

	//"time"
	//"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

//文章结构体
type Article struct {
	Title, Body string
	ID          int64
}

var router *mux.Router
var db *sql.DB

//验证表单内容函数

//博文存储部分函数

func main() {
	bootsrap.SetupDB()
	router = bootsrap.SetupRoute()
	// 中间件的使用 强转网页类型
	//通过命名路由获取URL（测试）
	// homeURL, _ := router.Get("home").URL()
	// fmt.Println("HomeURL:", homeURL)
	// articlesURL, _ := router.Get("articles.show").URL()
	// fmt.Println("ArticlesURL:", articlesURL)
	//URL去斜杠
	err := http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)
}
