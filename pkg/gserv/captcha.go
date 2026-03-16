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

type CaptchaTemplateGenerator struct {
	text string
	tmpl *template.Template
}

func NewCaptchaTemplateGenerator() *CaptchaTemplateGenerator {
	return &CaptchaTemplateGenerator{
		tmpl: template.New("captcha"),
	}
}

func (g *CaptchaTemplateGenerator) Open(path string) error {
	// 读取HTML模板
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	text, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	g.tmpl = template.Must(g.tmpl.Parse(string(text)))
	g.text = string(text)

	return nil
}

func (g *CaptchaTemplateGenerator) SendCaptchaEmail(to string, captcha string) error {
	var buf bytes.Buffer
	buf.WriteString(string(g.text))

	buf.Reset()

	err := g.tmpl.Execute(&buf, map[string]string{"Captcha": captcha})
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
