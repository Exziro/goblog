package user

import (
	"goblog/app/models"
	"goblog/pkg/logger"
	"goblog/pkg/types"
)

//Creat 创建用户，通过User.ID来判断是否创建成功
func (user *User) Creat() (err error) {
	if err = models.DB.Create(&user).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

// Get 通过ID获取用户
func Get(idstr string) (User, error) {
	var user User
	id := types.StringToUint64(idstr)
	if err := models.DB.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil

}

//GetByEmail 通过邮箱获取用户信息
func GetByEmail(email string) (User, error) {
	var user User
	if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// All 获取所有用户数据
func All() ([]User, error) {
	var users []User
	if err := models.DB.Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}
