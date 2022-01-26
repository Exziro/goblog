package request

import (
	"goblog/app/models/user"

	"github.com/thedevsaddam/govalidator"
)

//// ValidateRegistrationForm 验证表单，返回 errs 长度等于零即通过
func ValidateRegistrationForm(data user.User) map[string][]string {
	//定制认证规则
	rules := govalidator.MapData{
		"name":             []string{"required", "alpha_num", "between:3,20"},
		"email":            []string{"required", "min:4", "max:30", "email"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
	}
	//定制错误信息
	messages := govalidator.MapData{
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:格式错误，只允许数字和英文",
			"between:用户名长度需在 3~20 之间",
		},
		"password": []string{
			"required:密码为必填项",
			"min:长度需大于 6",
		},
		"password_confirm": []string{
			"required:确认密码框为必填项",
		},
	}
	//配置初始化
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid", //struct标签标识符
		Messages:      messages,
	}
	//开始验证
	errs := govalidator.New(opts).ValidateStruct()
	//因 govalidator 不支持 password_confirm 验证，我们自己写一个
	if data.Password != data.PasswordConfirm {
		errs["password_confirm"] = append(errs["password_confirm"], "两次输入密码不匹配")
	}
	return errs
}
