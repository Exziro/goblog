package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

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
func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章ID："+id)
}
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "文章列表")
}
func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "please check your form")
		return
	}
	//title := r.PostForm.Get("title")
	fmt.Fprintf(w, "r.Form 中 title 的值为: %v <br>", r.FormValue("title"))
	fmt.Fprintf(w, "r.PostForm 中 title 的值为: %v <br>", r.PostFormValue("title"))
	fmt.Fprintf(w, "r.Form 中 test 的值为: %v <br>", r.FormValue("test"))
	fmt.Fprintf(w, "r.PostForm 中 test 的值为: %v <br>", r.PostFormValue("test"))
}
func articlesCreatHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
    <title>创建文章 —— 我的技术博客</title>
	</head>
	<body>
    <form action="%s?test=data" method="post">
        <p><input type="text" name="title"></p>
        <p><textarea name="body" cols="30" rows="10"></textarea></p>
        <p><button type="submit">提交</button></p>
    </form>
	</body>
	</html>
	`
	storeURL, _ := router.Get("articles.store").URL()
	fmt.Fprintf(w, html, storeURL)
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
func main() {

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreatHandler).Methods("GET").Name("aricles.creat")
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
