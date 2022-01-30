package controllers

import (
	"fmt"
	"goblog/app/models/category"
	"goblog/app/request"
	"goblog/pkg/flash"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

//CategoriesController 文章分类控制器
type CategoriesController struct {
	BaseController
}

//Creat 文章分类创建页面
func (*CategoriesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "categories.create")
}

//Store 保存文章分类
func (*CategoriesController) Store(w http.ResponseWriter, r *http.Request) {
	//初始化数据
	_category := category.Category{
		Name: r.PostFormValue("name"),
	}
	//表单验证
	errors := request.ValidateCategoryForm(_category)
	//检测错误
	if len(errors) == 0 {
		//创建文章分类
		_category.Create()
		if _category.ID > 0 {
			flash.Success("创建分类成功")
			indexURL := route.Name2URL("home")
			http.Redirect(w, r, indexURL, http.StatusFound)
			// indexURL := route.Name2URL("categories.show", "id", _category.GetStringID())
			// http.Redirect(w, r, indexURL, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章分类失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Catrgory": _category,
			"Errors":   errors,
		}, "categories.create")
	}
}

//Show 显示分类栏
func (*CategoriesController) Show(w http.ResponseWriter, r *http.Request) {

}
