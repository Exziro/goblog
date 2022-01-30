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
