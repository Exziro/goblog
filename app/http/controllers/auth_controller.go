package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/pkg/view"
	"net/http"
)

// AuthController 处理静态页面
type AuthController struct {
}

//Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "auth.register")
}

//DoRegister 处理注册逻辑
func (*AuthController) DoRegiter(w http.ResponseWriter, r *http.Request) {
	//初始化变量
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("paasword")
	//表单验证
	//验证通过将数据存入
	_user := user.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	_user.Creat()
	//判定
	if _user.ID > 0 {
		fmt.Fprint(w, "插入成功")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "创建用户失败")
	}
	//表单不通过重新显示表单
}
