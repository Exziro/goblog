package main

import (
	"database/sql"
	//"go/types"
	//"goblog/pkg/logger"
	bootsrap "goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"

	//"goblog/pkg/model"

	"strconv"

	//"errors"
	"fmt"
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

//通过传参 id 获取博文
func getArticleByID(id string) (Article, error) {
	article := Article{}
	querry := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(querry, id).Scan(&article.ID, &article.Body, &article.Title)
	return article, err
}

//验证表单内容函数

func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.FormatInt(a.ID, 10))
	if err != nil {
		logger.LogError(err)
	}
	//删除成功，跳转到文章详情页
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}
	return 0, nil
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	//获取URL参数
	id := getRouterVariable("id", r)
	//读取对应的文章数据
	article, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 没有找到数据")
		} else {
			//
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}

	} else {
		//未出现错误，执行删除
		rowsAffected, err := article.Delete()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500服务器内部错误")

		} else {
			//未发生错误
			if rowsAffected > 0 {
				//重定向到文章列表
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)

			} else {
				//Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404未找到文章")

			}
		}
	}
}

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

func getRouterVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func main() {
	database.Initialize()
	db = database.DB
	bootsrap.SetupDB()
	router = bootsrap.SetupRoute()

	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")
	// 中间件的使用 强转网页类型
	router.Use(forceHTMLMiddleware)
	//通过命名路由获取URL（测试）
	// homeURL, _ := router.Get("home").URL()
	// fmt.Println("HomeURL:", homeURL)
	// articlesURL, _ := router.Get("articles.show").URL()
	// fmt.Println("ArticlesURL:", articlesURL)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
