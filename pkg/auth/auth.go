package auth

import (
	"errors"
	"goblog/app/models/user"
	"goblog/pkg/session"

	"gorm.io/gorm"
)

//获取uid
func _getUid() string {
	_uid := session.Get("uid")
	uid, ok := _uid.(string)
	if ok && len(uid) > 0 {
		return uid
	}
	return ""
}

//User 获取用户登录信息
func User() user.User {
	uid := _getUid()
	if len(uid) > 0 {
		_user, err := user.Get(uid)
		if err == nil {
			return _user
		}
	}
	return user.User{}
}

//Attempt 尝试登录
func Attempt(email string, password string) error {
	// 1. 根据 Email 获取用户
	_user, err := user.GetByEmail(email)

	// 2. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("账号不存在或密码错误")
		} else {
			return errors.New("内部错误，请稍后尝试")
		}
	}

	// 3. 匹配密码
	if !_user.ComparePassword(password) {
		return errors.New("账号不存在或密码错误")
	}

	// 4. 登录用户，保存会话
	session.Put("uid", _user.GetStringID())

	return nil
}

//Login 登录指定用户
func Login(_user user.User) {
	session.Put("uid", _user.GetStringID())

}

//Logout 退出用户
func Logout() {
	session.Forget("uid")
}

//Check 检测是否是登录状态
func Check() bool {
	return len(_getUid()) > 0
}
