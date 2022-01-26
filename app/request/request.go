package request

import (
	"errors"
	"fmt"
	"goblog/app/models"
	"strings"

	"github.com/thedevsaddam/govalidator"
)

func init() {
	// not_exists:users,email
	govalidator.AddCustomRule("not_exists", func(field string, rule string, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")

		tabelname := rng[0]
		dbField := rng[1]
		val := value.(string)

		var count int64
		models.DB.Table(tabelname).Where(dbField+" = ?", val).Count(&count)

		if count != 0 {

			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("%v 已被占用", val)
		}
		return nil
	})
}
