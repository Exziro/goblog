package category

import "goblog/app/models"

//文章分类模型
type Category struct {
	models.BaseModel

	Name string `gorm:"type:varchar(255);not null;" valid:"name"`
}
