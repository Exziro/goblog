package article

import (
	//"goblog/pkg/logger"
	"goblog/app/models"
	"goblog/pkg/route"
	"goblog/pkg/types"
)

//article文章类型
type Article struct {
	models.BaseModel
	Title string
	Body  string
}

//Link 模板中生成链接
func (article Article) Link() string {
	return route.Name2URL("articles.show", "id", article.GetStringID())
}
func (article Article) GetStringID() string {
	return types.Uint64ToString(article.ID)
}
