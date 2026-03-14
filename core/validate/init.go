package validate

import (
	"gServ/pkg/model"

	"github.com/go-playground/validator"
)

var (
	data_validator = validator.New()
)

func Init() error {
	if err := data_validator.RegisterValidation("captcha_type", validateTagCaptchaType); err != nil {
		return err
	}

	return nil
}

func validateTagCaptchaType(fl validator.FieldLevel) bool {
	captcha_type_value := fl.Field().Uint()

	switch captcha_type_value {
	case uint64(model.CAPTCHA_TYPE_UNKNOWN),
		uint64(model.CAPTCHA_TYPE_REGISTER),
		uint64(model.CAPTCHA_TYPE_RESET_PASSWORD),
		uint64(model.CAPTCHA_TYPE_CHANGE_EMAIL):
		return true
	}

	return false
}
