package model

import "time"

type CaptchaType uint8

const (
	CAPTCHA_TYPE_UNKNOWN CaptchaType = iota
	CAPTCHA_TYPE_REGISTER
	CAPTCHA_TYPE_RESET_PASSWORD
	CAPTCHA_TYPE_CHANGE_EMAIL
)

type EmailCaptcha struct {
	ModelHeader

	Email  string      `gorm:"type:varchar(255);not null" validate:"required,max=255,email"`
	Type   CaptchaType `validate:"captcha_type"`
	Code   string      `gorm:"type:varchar(255);not null" validate:"required,min=4,max=255"`
	UsedAt *time.Time

	ModelTail
}
