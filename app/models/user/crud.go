package user

import (
	"goblog/app/models"
	"goblog/pkg/logger"
)

//Creat 创建用户，通过User.ID来判断是否创建成功
func (user *User) Creat() (err error) {
	if err = models.DB.Create(&user).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}
