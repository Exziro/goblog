package view

import (
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

//Render 渲染视图
func Render(w io.Writer, name string, data interface{}) {
	//设置模板相对路径
	viewDir := "resources/views/"
	//2. 语法糖 将article.show更正为 articles/show
	name = strings.Replace(name, ".", "/", -1)
	//3 所有布局模板文件 Slice
	files, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)
	//4 在Slice里新增文件
	newfiles := append(files, viewDir+name+".gohtml")
	//5 解析所有模板
	tmpl, err := template.New(name + ".gohtml").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(newfiles...)
	logger.LogError(err)
	err = tmpl.ExecuteTemplate(w, "app", data)
	logger.LogError(err)
}
