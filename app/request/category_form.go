package request

import (
	"goblog/app/models/category"

	"github.com/thedevsaddam/govalidator"
)

// ValidateCategoryForm 验证表单，返回 errs 长度等于零即通过
func ValidateCategoryForm(data category.Category) map[string][]string {
	//定制认证规则
	rules := govalidator.MapData{
		"name": []string{"required", "min_cn:2", "max_cn:8", "not_exists:categories,name"},
	}
	//定制错误信息
	messages := govalidator.MapData{
		"name": []string{
			"required:分类名称为必填项",
			"min:分类名称长度需至少 2 个字",
			"max:分类名称长度不能超过 8 个字",
		},
	}
	//配置初始化
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid", // 模型中的 Struct 标签标识符

	}

	//返回验证
	return govalidator.New(opts).ValidateStruct()
}
