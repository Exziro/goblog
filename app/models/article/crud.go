package article

import (
	"goblog/app/models"
	"goblog/pkg/logger"

	"goblog/pkg/types"
)

//GET通过ID获取文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToUint64(idstr)
	if err := models.DB.First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}
func GetAll() ([]Article, error) {
	var articles []Article
	if err := models.DB.Find(&articles).Error; err != nil {
		return articles, err

	}
	return articles, nil
}
func (articles *Article) Create() (err error) {
	result := models.DB.Create(&articles)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}
func (articles *Article) Update() (rowsAffected int64, err error) {
	result := models.DB.Save(&articles)
	if err = result.Error; err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}
func (articles *Article) Delete() (rowsAffected int64, err error) {
	result := models.DB.Delete(&articles)
	if err = result.Error; err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}
