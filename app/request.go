package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/skszcool/iot-device/logger"
	"github.com/skszcool/iot-device/validation"
)

func MarkErrors(c *gin.Context, err error) string {
	var errStr string

	switch err.(type) {
	case validator.ValidationErrors:
		errStr = validation.Translate(err.(validator.ValidationErrors))
	case *json.UnmarshalTypeError:
		unmarshalTypeError := err.(*json.UnmarshalTypeError)
		errStr = fmt.Errorf("%s 类型错误，期望类型 %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String()).Error()
	default:
		errStr = errors.New("unknown error.").Error()
	}

	logger.Error(err)

	ResponseCheckParamsFailed(c, errStr, errStr)

	return errStr
}
