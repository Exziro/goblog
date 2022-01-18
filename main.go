package main

import (
	"database/sql"
	"goblog/pkg/route"
	"strconv"

	//"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog</h1>")

}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/html; charset= utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

//通过传参 URL 路由参数名称获取值
func getRouterVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

//通过传参 id 获取博文
func getArticleByID(id string) (Article, error) {
	article := Article{}
	querry := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(querry, id).Scan(&article.ID, &article.Body, &article.Title)
	return article, err
}

//RoutName2URL 通过路由名称获取URL
func RouteName2URL(routName string, pairs ...string) string {
	url, err := router.Get(routName).URL(pairs...)
	if err != nil {
		checkError(err)
		return ""
	}

	return url.String()
}

//int 64 转换为string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
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
func (a Article) Link() string {
	showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID, 10))
	if err != nil {
		checkError(err)
		return ""
	}
	return showURL.String()

}
func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.FormatInt(a.ID, 10))
	if err != nil {
		checkError(err)
	}
	//删除成功，跳转到文章详情页
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}
	return 0, nil
}
func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	//获取「URL」请求参数
	id := getRouterVariable("id", r)
	//读取对应文章数据
	article, err := getArticleByID(id)
	// article := Article{}
	// querry := "SELECT * FROM articles WHERE id = ?"
	// err := db.QueryRow(querry, id).Scan(&article.ID, &article.Body, &article.Title)
	//错误检测

	if err != nil {
		//未找到查找项目
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "数据未找到404")
		} else {
			//服务器内部错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器内部错误")
		}
	} else {
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL": route.Name2URL,
				"Int64ToString": Int64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)
		err = tmpl.Execute(w, article)
		checkError(err)
		//fmt.Fprint(w, "读取文章成功"+article.Title)
	}
	fmt.Fprint(w, "文章ID："+id)

}
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	//从数据库读取条目
	rows, err := db.Query("SELECT * from articles")
	checkError(err)
	//创建一个文章数组
	var articles []Article
	// 2.1 扫描每一行的结果并赋值到一个 article 对象中
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		checkError(err)
		//将新的内容追加进数组
		articles = append(articles, article)
	}
	//检测遍历时是否发生错误
	err = rows.Err()
	checkError(err)
	//加载模板
	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	checkError(err)
	//渲染模板
	err = tmpl.Execute(w, articles)
	checkError(err)
	fmt.Fprint(w, "文章列表")
}
func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	//表单验证
	errors := validateArticleFromData(title, body)
	if len(errors) == 0 {
		lastInsertID, err := saveArticlesToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprintf(w, "插入成功，ID为："+strconv.FormatInt(lastInsertID, 10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500服务器内部错误")
		}
		fmt.Fprint(w, "验证通过！")
		fmt.Fprintf(w, "title 的值为: %v <br>", title)
		fmt.Fprintf(w, "title 的长度为: %v <br>", utf8.RuneCountInString(title))
		fmt.Fprintf(w, "body 的值为: %v <br>", body)
		fmt.Fprintf(w, "body 的长度为: %v <br>", utf8.RuneCountInString(body))

	} else {
		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}

		// fmt.Fprintf(w, "有错误发生，errors 的值为: %v <br>", errors)
	}
	// err := r.ParseForm()
	// if err != nil {
	// 	fmt.Fprint(w, "please check your form")
	// 	return
	// }

	//title := r.PostForm.Get("title")
	//获取请求的数据由以下
	// fmt.Fprintf(w, "r.Form 中 title 的值为: %v <br>", r.FormValue("title"))
	// fmt.Fprintf(w, "r.PostForm 中 title 的值为: %v <br>", r.PostFormValue("title"))
	// fmt.Fprintf(w, "r.Form 中 test 的值为: %v <br>", r.FormValue("test"))
	// fmt.Fprintf(w, "r.PostForm 中 test 的值为: %v <br>", r.PostFormValue("test"))
}

//文章创建路由
func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
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
		checkError(err)
		err = tmp.Execute(w, data)
		checkError(err)
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
				checkError(err)
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
			checkError(err)
			err = tmpl.Execute(w, data)
			checkError(err)
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
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}

	} else {
		//未出现错误，执行删除
		rowsAffected, err := article.Delete()
		if err != nil {
			checkError(err)
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

//数据库链接初始化
func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "qiyixi19961016",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}
	//准备数据库连接池
	db, err = sql.Open("mysql", config.FormatDSN())
	//fmt.Printf(config.FormatDSN())
	checkError(err)
	//设置最大连接数
	db.SetMaxOpenConns(25)
	//设置最大空闲连接数
	db.SetMaxIdleConns(25)
	//设置每个链接的过期时间
	db.SetConnMaxIdleTime(5 * time.Minute)
	//
	err = db.Ping()
	checkError(err)
}

//建表函数
func createTables() {
	//var err error
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
		id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
		title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
		body longtext COLLATE utf8mb4_unicode_ci
	); `
	_, err := db.Exec(createArticlesSQL)
	checkError(err)
}

//博文存储部分函数
func saveArticlesToDB(title, body string) (int64, error) {
	//变量初始化
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)
	//获取一个prepare
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
	//例行错误检查
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}
	return 0, nil
}

//报错函数
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {

	initDB()
	createTables()
	route.Initialize()
	router = route.Router
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("aricles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesHandlerEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")
	//自定义404
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	// 中间件的使用 强转网页类型
	router.Use(forceHTMLMiddleware)
	//通过命名路由获取URL（测试）
	// homeURL, _ := router.Get("home").URL()
	// fmt.Println("HomeURL:", homeURL)
	// articlesURL, _ := router.Get("articles.show").URL()
	// fmt.Println("ArticlesURL:", articlesURL)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
