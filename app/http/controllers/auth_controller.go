package controllers

import (
	"fmt"

	"goblog/app/models/user"
	"goblog/app/request"
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"goblog/pkg/session"
	"goblog/pkg/view"
	"net/http"
)

// AuthController 处理静态页面
type AuthController struct {
}

// type userForm struct {
// 	Name            string `valid:"name"`
// 	Email           string `valid:"email"`
// 	Password        string `valid:"password"`
// 	PasswordConfirm string `valid:"password_confirm"`
// }

//Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

//DoRegister 处理注册逻辑
func (*AuthController) DoRegiter(w http.ResponseWriter, r *http.Request) {
	//初始化变量
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}

	//表单验证加规则设定验证通过将数据存入
	errs := request.ValidateRegistrationForm(_user)
	if len(errs) > 0 {
		//错误发生 打印错误
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		//验证成功，创建数据
		_user.Creat()
		if _user.ID > 0 {
			flash.Success("恭喜你注册成功！")
			auth.Login(_user)
			//fmt.Fprint(w, "插入成功")

			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建用户失败")
		}
	}
}

//Login 现实用户登录表单
func (au *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	session.Flush()
	view.RenderSimple(w, view.D{}, "auth.login")
}

//DoLogin 验证用户登录
func (au *AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	//初始化数据表单

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	//尝试登录

	if err := auth.Attempt(email, password); err == nil {
		flash.Success("欢迎回来")
		//登陆成功 跳转
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		//登录失败，显示错误表单
		view.RenderSimple(w, view.D{
			"Error":    err.Error(),
			"Email":    email,
			"Password": password,
		}, "auth.login")
	}

}

//Logout登出操作
func (au *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout()
	flash.Success("成功退出")
	http.Redirect(w, r, "/", http.StatusFound)
}
