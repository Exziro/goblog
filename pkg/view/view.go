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
func Render(w io.Writer, data interface{}, tplFiles ...string) {
	//设置模板相对路径
	viewDir := "resources/views/"
	//2. 语法糖 将article.show更正为 articles/show
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}
	//3 所有布局模板文件 Slice
	layoutFiles, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)
	//4 在Slice里新增文件 合并所有文件
	allFiles := append(layoutFiles, tplFiles...)
	//5 解析所有模板
	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(allFiles...)
	logger.LogError(err)
	err = tmpl.ExecuteTemplate(w, "app", data)
	logger.LogError(err)
}
