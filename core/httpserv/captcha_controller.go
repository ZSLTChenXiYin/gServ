package httpserv

import (
	"gServ/core/repository"
	"gServ/core/validate"
	"gServ/pkg/gserv"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 发送邮箱验证码
func post_Api_Captcha_Email(c *gin.Context) {
	req := &post_Api_Captcha_Email_Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code := captcha_generator.Generate()
	err := repository.CreateEmailCaptcha(req.Email, req.EmailType, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成验证码失败"})
		return
	}

	err = gserv.SendCaptchaEmail(req.Email, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送验证码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送"})
}
