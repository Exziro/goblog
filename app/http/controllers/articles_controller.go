package controllers

import (
	//"database/sql"
	"fmt"
	"goblog/pkg/logger"
	"path/filepath"

	"strconv"
	"unicode/utf8"

	"goblog/pkg/model/article"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"net/http"
	"text/template"

	"gorm.io/gorm"
)

type ArticlesController struct {
}

//ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         string
	Errors      map[string]string
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
		//设置模板相对路径
		viewDir := "resources/views"

		//所有布局模板文件slice
		files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
		logger.LogError(err)
		//在 Slice 里新增我们的目标文件
		newFiles := append(files, viewDir+"/articles/index.gohtml")
		//解析模板文件
		tmpl, err := template.ParseFiles(newFiles...)
		logger.LogError(err)

		//渲染模板
		err = tmpl.ExecuteTemplate(w, "app", articles)
		logger.LogError(err)
	}
}

//Store文章保存
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	//表单验证
	errors := validateArticleFromData(title, body)
	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprintf(w, "插入成功，ID为："+strconv.FormatUint(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500服务器内部错误")
		}
		fmt.Fprint(w, "验证通过！")
		fmt.Fprintf(w, "title 的值为: %v <br>", title)
		fmt.Fprintf(w, "title 的长度为: %v <br>", utf8.RuneCountInString(title))
		fmt.Fprintf(w, "body 的值为: %v <br>", body)
		fmt.Fprintf(w, "body 的长度为: %v <br>", utf8.RuneCountInString(body))

	} else {
		storeURL := route.Name2URL("articles.store")

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		logger.LogError(err)

		err = tmpl.Execute(w, data)
		logger.LogError(err)

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

//Create 文章创建
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	storeURL := route.Name2URL("articles.store")
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

//Edit 博文修改
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]
	id := route.GetRouterVariable("id", r)
	//读取文章数据
	article, err := article.Get(id)
	// article := Article{}
	// querry := "SELECT * FROM articles WHERE id = ?"
	// err := db.QueryRow(querry, id).Scan(&article.ID, &article.Title, &article.Body)

	//查错
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "文章未查找到404")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器内部错误500")
		}

	} else {
		// 4. 读取成功，显示编辑文章表单
		updateURL := route.Name2URL("articles.update", "id", id)
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

//更新博文
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouterVariable("id", r)
	_article, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
			_article.Title = title
			_article.Body = body

			rowsAffected, err := _article.Update()
			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}
			//更新成功进入跳转页面
			if rowsAffected > 0 {
				showURL := route.Name2URL("articles.show", "id", id)
				http.Redirect(w, r, showURL, http.StatusFound)

			} else {
				fmt.Fprint(w, "没有做任何更改")
			}

		} else {
			//表单验证不通过，显示理由
			updateURL := route.Name2URL("articles.update", "id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				Errors: errors,
				URL:    updateURL,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)
			err = tmpl.Execute(w, data)
			logger.LogError(err)
			fmt.Fprintf(w, "更新失败")
		}
	}

}

//Delete 删除博文
func (articles *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	//获取URL参数
	id := route.GetRouterVariable("id", r)
	//读取对应的文章数据
	_article, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//数据未找到
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
		rowsAffected, err := _article.Delete()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500服务器内部错误")

		} else {
			//未发生错误
			if rowsAffected > 0 {
				//重定向到文章列表
				indexURL := route.Name2URL("articles.index", "id", id)
				http.Redirect(w, r, indexURL, http.StatusFound)

			} else {
				//Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404未找到文章")

			}
		}
	}
}
