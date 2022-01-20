package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RouteName2URL 通过路由名称来获取 URL
func Name2URL(routName string, pairs ...string) string {
	var route *mux.Router
	url, err := route.Get(routName).URL(pairs...)
	if err != nil {
		//checkError(err)
		return ""
	}

	return url.String()
}

//通过传参 URL 路由参数名称获取值
func GetRouterVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
