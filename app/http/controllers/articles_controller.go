package controllers

import (
	//"database/sql"
	"fmt"
	"goblog/pkg/logger"

	"goblog/pkg/model/article"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"net/http"
	"text/template"

	"gorm.io/gorm"
)

type ArticlesController struct {
}

//展示文章（博文）内容
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	//获取「URL」请求参数
	id := route.GetRouterVariable("id", r)
	//读取对应文章数据
	article, err := article.Get(id)
	// article := Article{}
	// querry := "SELECT * FROM articles WHERE id = ?"
	// err := db.QueryRow(querry, id).Scan(&article.ID, &article.Body, &article.Title)
	//错误检测

	if err != nil {
		//未找到查找项目
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "数据未找到404")
		} else {
			//服务器内部错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器内部错误500")
		}
	} else {
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL":  route.Name2URL,
				"Uint64ToString": types.Uint64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)
		err = tmpl.Execute(w, article)
		logger.LogError(err)
		//fmt.Fprint(w, "读取文章成功"+article.Title)
	}
	//fmt.Fprint(w, "文章ID："+id)

}

//Index文章列表
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()
	if err != nil {
		//数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500服务器错误")
	} else {
		//加载模板
		tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		logger.LogError(err)

		//渲染模板
		err = tmpl.Execute(w, articles)
		logger.LogError(err)
	}
}
