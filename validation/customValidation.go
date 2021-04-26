// 自定义参数校验方法
package validation

import (
	"github.com/go-playground/validator/v10"
	"strconv"
)

func checkMobile(fl validator.FieldLevel) bool {
	mobile := strconv.Itoa(int(fl.Field().Uint()))
	if len(mobile) != 11 {
		return false
	}
	return true
}
