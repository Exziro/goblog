package main

import (
	"database/sql"
	//"go/types"
	//"goblog/pkg/logger"
	bootsrap "goblog/bootstrap"
	"goblog/pkg/database"

	//"goblog/pkg/model"

	//"errors"

	//"html/template"
	"net/http"
	"net/url"
	"strings"

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

//中间件处理，用于设置所有页面适配请求头的处理模式
func forceHTMLMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//设置标头
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		//继续处理请求
		h.ServeHTTP(w, r)
	})
}

// 把 Gorilla Mux 包起来，在这个函数中我们先对进来的请求做处理，然后再传给 Gorilla Mux 去解析
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	})
}

//博文存储部分函数

func main() {
	database.Initialize()
	db = database.DB
	bootsrap.SetupDB()
	router = bootsrap.SetupRoute()

	// 中间件的使用 强转网页类型
	router.Use(forceHTMLMiddleware)
	//通过命名路由获取URL（测试）
	// homeURL, _ := router.Get("home").URL()
	// fmt.Println("HomeURL:", homeURL)
	// articlesURL, _ := router.Get("articles.show").URL()
	// fmt.Println("ArticlesURL:", articlesURL)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
