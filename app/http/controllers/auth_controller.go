package controllers

import (
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

}
