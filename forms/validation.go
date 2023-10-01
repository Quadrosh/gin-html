package forms

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() *validator.Validate {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("max_length", MaxLength)
		return v
	}
	return nil
}

var MaxLength validator.Func = func(fl validator.FieldLevel) bool {
	str, ok := fl.Field().Interface().(string)
	if ok {
		var paramStr = fl.Param()
		uintPar, err := strconv.Atoi(paramStr)
		if err != nil {
			log.Println(err.Error())
			return false
		}
		var rStr = []rune(str)
		if len(rStr) > uintPar {
			return false
		}
	}
	return true
}
