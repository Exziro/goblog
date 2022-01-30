package route

import (
	"goblog/pkg/config"
	"goblog/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

var route *mux.Router

func SetRoute(r *mux.Router) {
	route = r
}

// RouteName2URL 通过路由名称来获取 URL
func Name2URL(routName string, pairs ...string) string {

	url, err := route.Get(routName).URL(pairs...)
	if err != nil {
		//checkError(err)
		logger.LogError(err)
		return ""
	}

	return config.GetString("app.url") + url.String()
}

//通过传参 URL 路由参数名称获取值
func GetRouterVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
