package gserv

import (
	"bytes"
	"fmt"
	"gServ/core/config"
	"gServ/core/log"
	"io"
	"math/rand"
	"os"
	"text/template"
	"time"

	"gopkg.in/gomail.v2"
)

const (
	CAPTCHA_LENGTH = 4
	CAPTCHA_DIC    = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

type CaptchaGenerator struct {
	length uint
	rander *rand.Rand
}

func NewCaptchaGenerator(length uint) *CaptchaGenerator {
	return &CaptchaGenerator{
		length: length,
		rander: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *CaptchaGenerator) Generate() string {
	captcha := make([]byte, g.length)
	for i := 0; i < int(g.length); i++ {
		captcha[i] = CAPTCHA_DIC[g.rander.Intn(len(CAPTCHA_DIC))]
	}
	return string(captcha)
}

func SendCaptchaEmail(to string, captcha string) error {
	// 读取HTML模板
	file, err := os.Open("captcha.html")
	if err != nil {
		return err
	}
	tmpl, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	t := template.Must(template.New("email-captcha").Parse(string(tmpl)))

	var buf bytes.Buffer
	buf.WriteString(string(tmpl))

	buf.Reset()

	err = t.Execute(&buf, map[string]string{"Captcha": captcha})
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("浊水楼台-gServ <%s>", config.GetConfig().Server.Email.Email))
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "浊水楼台-gServ 验证码")
	msg.SetBody("text/html", buf.String())

	dialer := gomail.NewDialer(
		config.GetConfig().Server.Email.Host,
		int(config.GetConfig().Server.Email.Port),
		config.GetConfig().Server.Email.Email,
		config.GetConfig().Server.Email.Password,
	)

	go func() {
		err := dialer.DialAndSend(msg)
		if err != nil {
			log.StdErrorf("邮件发送失败: %v", err)
		}
	}()

	return nil
}
