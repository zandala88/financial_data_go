package vaildator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func DateValidator(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	// 正则表达式验证日期格式 YYYY-MM-DD
	regex := `^\d{4}-\d{2}-\d{2}$`
	matched, _ := regexp.MatchString(regex, date)
	return matched
}
