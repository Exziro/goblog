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
	"html/template"
	"net/http"
	"net/url"
	"strings"

	//"time"
	"unicode/utf8"

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
func validateArticleFromData(title, body string) map[string]string {
	errors := make(map[string]string)
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString("title") < 3 || utf8.RuneCountInString("title") > 40 {
		errors["title"] = "标题字数不正确"
	}
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}
	return errors
}

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

//博文修改
func articlesHandlerEditHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]
	id := getRouterVariable("id", r)
	//读取文章数据
	article, err := getArticleByID(id)
	// article := Article{}
	// querry := "SELECT * FROM articles WHERE id = ?"
	// err := db.QueryRow(querry, id).Scan(&article.ID, &article.Title, &article.Body)

	//查错
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "文章未查找到404")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器内部错误500")
		}

	} else {
		updateURL, _ := router.Get("articles.update").URL("id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}
		tmp, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)
		err = tmp.Execute(w, data)
		logger.LogError(err)
	}
}
func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouterVariable("id", r)
	_, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "300 服务器内部出错")
		}

	} else {
		title := r.FormValue("title")
		body := r.FormValue("body")
		//表单验证
		errors := validateArticleFromData(title, body)

		if len(errors) == 0 {
			//表单验证结束 将内容进行更新
			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
			rs, err := db.Exec(query, title, body, id)
			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}
			//更新成功进入跳转页面
			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)

			} else {
				fmt.Fprint(w, "没有做任何更改")
			}

		} else {
			//表单验证不通过，显示理由
			upDateURL, _ := router.Get("articles.update").URL()
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				Errors: errors,
				URL:    upDateURL,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)
			err = tmpl.Execute(w, data)
			logger.LogError(err)
			fmt.Fprintf(w, "更新失败")
		}
	}
	//fmt.Fprintf(w, "更新成功")
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

	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesHandlerEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
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
