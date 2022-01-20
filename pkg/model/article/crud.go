package article

import (
	//"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"
)

//GET通过ID获取文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToUint64(idstr)
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}
