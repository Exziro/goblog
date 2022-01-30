package article

import (
	//"goblog/pkg/logger"
	"goblog/app/models"
	"goblog/app/models/user"
	"goblog/pkg/route"
	"goblog/pkg/types"
)

//article文章类型
type Article struct {
	models.BaseModel
	Title      string `gorm:"type:varchar(255);not null;" valid:"title"`
	Body       string `gorm:"type:longtext;not null;" valid:"body"`
	UserID     uint64 `gorm:"not null;index"`
	CategoryID uint64 `gorm:"not null;default:4;index"`
	User       user.User
}

//Link 模板中生成链接
func (article Article) Link() string {
	return route.Name2URL("articles.show", "id", article.GetStringID())
}
func (article Article) GetStringID() string {
	return types.Uint64ToString(article.ID)
}

//CreatedAtDate 创建日期
func (article Article) CreatedAtDate() string {
	return article.CreatedAt.Format("2006-01-02")
}
