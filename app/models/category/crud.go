package category

import (
	"goblog/app/models"
	"goblog/pkg/logger"
)

// Create 创建分类，通过 category.ID 来判断是否创建成功
func (category *Category) Create() (err error) {
	if err = models.DB.Create(&category).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

// All 获取分类数据
func All() ([]Category, error) {
	var categories []Category
	if err := models.DB.Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}
