package repository

import (
	"gServ/pkg/model"
	"time"
)

func CreateEmailCaptcha(email string, captcha_type model.CaptchaType, code string) error {
	return database.Create(&model.EmailCaptcha{
		Email: email,
		Type:  captcha_type,
		Code:  code,
	}).Error
}

func FindEmailCaptchasByEmailAndCaptchaType(email string, captcha_type model.CaptchaType) ([]model.EmailCaptcha, error) {
	var captchas []model.EmailCaptcha
	// 判断是否使用过
	err := database.Where("email = ? AND type = ? AND used_at IS NULL", email, captcha_type).Find(&captchas).Error
	return captchas, err
}

func UpdateEmailCaptchaUsedAt(id uint) error {
	return database.Model(&model.EmailCaptcha{}).
		Where("id = ?", id).Update("used_at", time.Now()).Error
}

// 删除过期的邮箱验证码（创建五分钟后视为过期）
func DeleteUnusedEmailCaptchas() error {
	return database.Model(&model.EmailCaptcha{}).
		Where("created_at < ?", time.Now().Add(-5*time.Minute)).
		Delete(&model.EmailCaptcha{}).Error
}
