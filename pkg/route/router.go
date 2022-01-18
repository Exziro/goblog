package route

import "github.com/gorilla/mux"

var Router *mux.Router

func Initialize() {
	Router = mux.NewRouter()
}

// RouteName2URL 通过路由名称来获取 URL
func Name2URL(routName string, pairs ...string) string {
	url, err := Router.Get(routName).URL(pairs...)
	if err != nil {
		//checkError(err)
		return ""
	}

	return url.String()
}
